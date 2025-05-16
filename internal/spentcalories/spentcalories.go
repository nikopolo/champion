package spentcalories

import (
	"errors"
	"fmt"
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

	if data == "" {
		return 0, "", 0, errors.New("неверный формат данных")
	}
	arr := strings.Split(data, ",")
	if len(arr) != 3 {
		return 0, "", 0, fmt.Errorf("указанны не все параметры")
	}
	count, err := strconv.Atoi(arr[0]) // Количество шагов
	if err != nil {
		return 0, "", 0, err
	}

	activity := arr[1]                              // Активность
	walkDuration, err := time.ParseDuration(arr[2]) // Продолжительность прогулки
	if err != nil {
		return 0, "", 0, err
	}

	if count <= 0 || activity == "" || walkDuration <= 0 {
		return 0, "", 0, fmt.Errorf("неверные данные")
	}

	return count, activity, walkDuration, nil
}

func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient          // Длина шага
	distanceInKm := (float64(steps) * stepLen) / mInKm // Дистанция в километрах
	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)

	return distance / duration.Hours() // продолжительность в часах
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	countSteps, activity, walkDuration, err := parseTraining(data)

	if err != nil {
		return "", err
	}

	var (
		calcDistance  float64
		calcMeanSpeed float64
		calcCalories  float64
	)

	switch activity {
	case "Бег":
		calcDistance = distance(countSteps, height)
		calcMeanSpeed = meanSpeed(countSteps, height, walkDuration)
		calcCalories, err = RunningSpentCalories(countSteps, weight, height, walkDuration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		calcDistance = distance(countSteps, height)
		calcMeanSpeed = meanSpeed(countSteps, height, walkDuration)
		calcCalories, err = WalkingSpentCalories(countSteps, weight, height, walkDuration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, walkDuration.Hours(), calcDistance, calcMeanSpeed, calcCalories), nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные данные")
	}
	speed := meanSpeed(steps, height, duration) // Средняя скорость
	return (weight * speed * duration.Minutes()) / minInH, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные данные")
	}
	speed := meanSpeed(steps, height, duration) // Средняя скорость
	calories := (weight * speed * duration.Minutes()) / minInH
	return calories * walkingCaloriesCoefficient, nil
}
