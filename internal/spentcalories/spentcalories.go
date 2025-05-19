package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		log.Printf("error: invalid format for string (steps,activity,time): %q", data)
		return 0, "", 0, fmt.Errorf("invalid format: expected steps,activity,time")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Printf("error: invalid steps value: %q, %v", parts[0], err)
		return 0, "", 0, fmt.Errorf("invalid steps value: %w", err)
	}

	if steps <= 0 {
		log.Printf("error: steps must be positive: %d", steps)
		return 0, "", 0, fmt.Errorf("steps must be positive")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		log.Printf("error: invalid duration: %q, %v", parts[2], err)
		return 0, "", 0, fmt.Errorf("invalid duration: %w", err)
	}

	if duration == 0 {
		log.Printf("error: duration cannot be zero: %q", parts[2])
		return 0, "", 0, fmt.Errorf("duration cannot be zero")
	}

	if duration < 0 {
		log.Printf("error: duration cannot be negative: %q", parts[2])
		return 0, "", 0, fmt.Errorf("duration cannot be negative")
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	durationHours := duration.Hours()
	return dist / durationHours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, fmt.Errorf("steps must be positive")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("weight must be positive")
	}
	if height <= 0 {
		return 0, fmt.Errorf("height must be positive")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("duration must be positive")
	}

	meanSpeed := meanSpeed(steps, height, duration)

	durationMinutes := duration.Minutes()
	calories := (weight * meanSpeed * durationMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("steps must be positive")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("weight must be positive")
	}
	if height <= 0 {
		return 0, fmt.Errorf("height must be positive")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("duration must be positive")
	}

	meanSpeed := meanSpeed(steps, height, duration)

	durationMinutes := duration.Minutes()
	calories := (weight * meanSpeed * durationMinutes) / minInH * walkingCaloriesCoefficient
	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Printf("error: failed to parse training data: %v", err)
		return "", err
	}

	var calories float64
	var activityName string
	switch activityType {
	case "Ходьба":
		activityName = "Ходьба"
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Printf("error: failed to calculate walking calories: %v", err)
			return "", err
		}
	case "Бег":
		activityName = "Бег"
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Printf("error: failed to calculate running calories: %v", err)
			return "", err
		}
	default:
		log.Printf("error: unknown activity type: %q", activityType)
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	dist := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, duration)

	durationHours := duration.Hours()

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityName, durationHours, dist, meanSpeed, calories), nil
}
