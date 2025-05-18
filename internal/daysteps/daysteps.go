package daysteps

import (
	"errors"
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
		log.Println("Ошибка: неверный формат строки (шаги,время)")
		return 0, 0, errors.New("Неправильный формат: шаги, время")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Println("Ошибка: некорректное количество шагов")
		return 0, 0, errors.New("invalid steps: " + err.Error())
	}

	if steps <= 0 {
		log.Println("Ошибка: некорректное количество шагов")
		return 0, 0, errors.New("некорректное количество шагов")
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		log.Println("Ошибка: некорректная продолжительность тренировки")
		return 0, 0, errors.New("некорректная продолжительность: " + err.Error())
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 {
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
