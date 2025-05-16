package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	if data == "" {
		return 0, 0, errors.New("неверный тиg данных")
	}

	arr := strings.Split(data, ",")
	if len(arr) != 2 {
		return 0, 0, errors.New("неверный тиg данных")
	}

	countSteps, err := strconv.Atoi(arr[0]) // Количество шагов
	if err != nil {
		return 0, 0, err
	}

	if countSteps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	walkDuration, err := time.ParseDuration(arr[1]) // Продолжительность прогулки
	if err != nil {
		return 0, 0, err
	}
	if walkDuration <= 0 {
		return 0, 0, errors.New("неверная продолжительность")
	}

	return countSteps, walkDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	count, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return ""
	}
	if count <= 0 {
		log.Println(err)
		return ""
	}

	distance := (float64(count) * stepLength) / mInKm // Дистанция

	calories, err := spentcalories.WalkingSpentCalories(count, weight, height, duration) //Калории
	if err != nil {
		log.Println(err)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", count, distance, calories)
	return result
}
