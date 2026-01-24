package constants

// Общие ошибки
const (
	ErrGameNotFound        = "игра не найдена"
	ErrGameIDRequired      = "ID игры обязателен"
	ErrTitleRequired       = "название игры обязательно"
	ErrTitleTooLong        = "название игры слишком длинное (максимум 200 символов)"
	ErrDescriptionTooLong  = "описание игры слишком длинное (максимум 2000 символов)"
	ErrReleaseDateTooFar   = "дата релиза слишком далеко в будущем (максимум 1 год)"
	ErrIDMismatch          = "несоответствие ID игры"
	ErrGameAlreadyExists   = "игра уже существует"
	ErrGameCannotBeDeleted = "игра не может быть удалена"
	ErrInvalidData         = "неверные данные"
	ErrUnauthorized        = "неавторизованный доступ"
	ErrForbidden           = "доступ запрещен"

	ErrGenreNotFound   = "жанр не найден"
	ErrGenreIDRequired = "ID жанра обязателен"
	ErrGenreTitleEmpty = "название жанра обязательно"
)

// Ошибки валидации
const (
	ErrValidationTitleLength = "название должно содержать от 1 до 200 символов"
	ErrValidationDescription = "описание не должно превышать 2000 символов"
	ErrValidationReleaseDate = "дата релиза не может быть больше чем на год вперед"
	ErrValidationGenreID     = "ID жанра обязателен"
	ErrValidationLimit       = "лимит не может быть отрицательным"
	ErrValidationOffset      = "смещение не может быть отрицательным"
	ErrValidationMaxLimit    = "лимит не может превышать 100"
)

// Бизнес-ошибки
const (
	ErrBusinessCannotReleaseOnWeekend = "игры нельзя выпускать по выходным"
	ErrBusinessDailyLimitExceeded     = "превышен дневной лимит создания игр"
	ErrBusinessGameNotReleased        = "игра еще не выпущена"
	ErrBusinessDuplicateTitle         = "игра с таким названием уже существует"
)

// Сообщения успеха
const (
	MsgGameCreated      = "Игра успешно создана"
	MsgGameUpdated      = "Игра успешно обновлена"
	MsgGameDeleted      = "Игра успешно удалена"
	MsgOperationSuccess = "Операция выполнена успешно"
)
