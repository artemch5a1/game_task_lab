import { onCleanup, onMount, createEffect } from "solid-js";
import flatpickr from "flatpickr";
import type { Instance } from "flatpickr/dist/types/instance";
import "flatpickr/dist/flatpickr.min.css";

interface FlatpickrInputProps {
  value: string;
  onChange: (date: string) => void;
  disabled?: boolean;
  required?: boolean;
  placeholder?: string;
  id?: string;
  class?: string;
  style?: Record<string, string>;
}

export const FlatpickrInput = (props: FlatpickrInputProps) => {
  let inputRef: HTMLInputElement | undefined;
  let flatpickrInstance: Instance | null = null;

  onMount(() => {
    if (inputRef) {
      // Устанавливаем начальное значение в input перед инициализацией
      const initialValue = props.value || "";
      inputRef.value = initialValue;

      flatpickrInstance = flatpickr(inputRef, {
        dateFormat: "Y-m-d",
        defaultDate: initialValue || undefined,
        locale: {
          firstDayOfWeek: 1,
          weekdays: {
            shorthand: ["Вс", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб"],
            longhand: [
              "Воскресенье",
              "Понедельник",
              "Вторник",
              "Среда",
              "Четверг",
              "Пятница",
              "Суббота",
            ],
          },
          months: {
            shorthand: [
              "Янв",
              "Фев",
              "Мар",
              "Апр",
              "Май",
              "Июн",
              "Июл",
              "Авг",
              "Сен",
              "Окт",
              "Ноя",
              "Дек",
            ],
            longhand: [
              "Январь",
              "Февраль",
              "Март",
              "Апрель",
              "Май",
              "Июнь",
              "Июль",
              "Август",
              "Сентябрь",
              "Октябрь",
              "Ноябрь",
              "Декабрь",
            ],
          },
        },
        onChange: (_selectedDates, dateStr) => {
          if (dateStr) {
            props.onChange(dateStr);
          }
        },
        disableMobile: true,
      });

      // Принудительно устанавливаем значение после инициализации
      // Используем несколько попыток для надежности
      if (initialValue) {
        const setValue = () => {
          if (flatpickrInstance && inputRef) {
            flatpickrInstance.setDate(initialValue, false);
            // Убеждаемся, что значение отображается в input
            inputRef.value = initialValue;
          }
        };
        
        // Пробуем установить сразу
        setValue();
        
        // И еще раз после небольшой задержки
        setTimeout(setValue, 100);
        requestAnimationFrame(setValue);
      }
    }
  });

  // Отслеживаем изменения value извне
  createEffect(() => {
    const value = props.value;
    if (flatpickrInstance && inputRef) {
      // Используем небольшую задержку для надежности
      const timeoutId = setTimeout(() => {
        if (flatpickrInstance && inputRef) {
          const currentValue = flatpickrInstance.input.value;
          if (value && currentValue !== value) {
            flatpickrInstance.setDate(value, false);
            // Дублируем установку значения напрямую в input
            if (inputRef.value !== value) {
              inputRef.value = value;
            }
          } else if (!value && currentValue) {
            flatpickrInstance.clear();
            inputRef.value = "";
          }
        }
      }, 100);

      return () => clearTimeout(timeoutId);
    }
  });

  onCleanup(() => {
    if (flatpickrInstance) {
      flatpickrInstance.destroy();
      flatpickrInstance = null;
    }
  });

  return (
    <input
      ref={inputRef}
      type="text"
      id={props.id}
      class={props.class}
      style={{
        color: "#111827",
        ...props.style,
      }}
      placeholder={props.placeholder}
      required={props.required}
      disabled={props.disabled}
      readonly
    />
  );
};
