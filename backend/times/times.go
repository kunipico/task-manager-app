package times

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"task-manager-api/db"
	"task-manager-api/task"
)

type TimeInfo struct {
	Today     string `json:"today"`
	ThisWeek  string `json:"thisWeek"`
	ThisMonth string `json:"thisMonth"`
}

// calculateTaskTime は、指定された期間に基づいて時間を積算します
func calculateTaskTime(taskID int, startDate time.Time) (time.Duration, error) {
	query := `
		SELECT SetStatus, SetTime
		FROM Times
		WHERE Task_ID = ? AND SetTime >= ?
		ORDER BY SetTime ASC
	`
	rows, err := db.DB.Query(query, taskID, startDate)
	if err != nil {
		return 0, fmt.Errorf("データベースクエリエラー: %v", err)
	}
	defer rows.Close()

	var (
		totalDuration time.Duration
		prevStartTime time.Time
		inProgress    bool
	)

	for rows.Next() {
		var status string
		var SetTime time.Time

		if err := rows.Scan(&status, &SetTime); err != nil {
			return 0, fmt.Errorf("行の読み取りエラー: %v", err)
		}

		if status == "Start" {
			prevStartTime = SetTime
			inProgress = true
		} else if status == "Stop" && inProgress {
			totalDuration += SetTime.Sub(prevStartTime)
			inProgress = false
		}
	}

	// 処理が途中で終了した場合、現在時刻までを計算
	if inProgress {
		totalDuration += time.Since(prevStartTime)
	}

	return totalDuration, nil
}



func GetTimeInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Time Info")
	_, idStr, err := task.ExtractUserIDAndParam(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
		return
	}
	fmt.Println("idStr : ",idStr)

	// URLパスからパラメータを取得
	idpara := strings.TrimPrefix(idStr, "time-info/")
	
	fmt.Println("idpara : ",idpara)

	id, err := strconv.Atoi(idpara)
	fmt.Println("id : ",id)
	if err != nil {
		http.Error(w, "無効なIDです1", http.StatusBadRequest)
		return
	}

	now := time.Now()

	// 今日
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayDuration, err := calculateTaskTime(id, todayStart)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// 今週 (月曜日始まり)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 日曜日を週の最後にする
	}
	weekStart := todayStart.AddDate(0, 0, -weekday+1)
	weekDuration, err := calculateTaskTime(id, weekStart)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// 今月
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthDuration, err := calculateTaskTime(id, monthStart)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// JSON レスポンス
	timeInfo := TimeInfo{
		Today:    fmt.Sprintf("%d時間%d分", int(todayDuration.Hours()), int(todayDuration.Minutes())%60),
		ThisWeek: fmt.Sprintf("%d時間%d分", int(weekDuration.Hours()), int(weekDuration.Minutes())%60),
		ThisMonth: fmt.Sprintf("%d時間%d分", int(monthDuration.Hours()), int(monthDuration.Minutes())%60),
	}

	fmt.Println("timeInfo: ",timeInfo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(timeInfo)
}