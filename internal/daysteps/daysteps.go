package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"18_go_tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		log.Printf("error: invalid format for string (steps,time): %q", data)
		return 0, 0, fmt.Errorf("invalid format: expected steps,time")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Printf("error: invalid steps value: %q, %v", parts[0], err)
		return 0, 0, fmt.Errorf("invalid steps value: %w", err)
	}

	if steps <= 0 {
		log.Printf("error: steps must be positive: %d", steps)
		return 0, 0, fmt.Errorf("steps must be positive")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		log.Printf("error: invalid duration: %q, %v", parts[1], err)
		return 0, 0, fmt.Errorf("invalid duration: %w", err)
	}
	if duration == 0 {
		log.Printf("error: duration cannot be zero: %q", parts[1])
		return 0, 0, fmt.Errorf("duration cannot be zero")
	}

	if duration < 0 {
		log.Printf("error: duration cannot be negative: %q", parts[1])
		return 0, 0, fmt.Errorf("duration cannot be negative")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distanceMeters := float64(steps) * stepLength

	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)
}
