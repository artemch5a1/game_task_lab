import type { JSX } from "solid-js";
import { For, Show } from "solid-js";
import "./Table.css";

export interface TableColumn {
  key: string;
  header: string;
  align?: "left" | "center" | "right";
  render?: (row: any) => JSX.Element;
}

interface TableProps {
  columns: TableColumn[];
  data: any[];
  isLoading?: boolean;
  emptyText?: string;
}

export const Table = (props: TableProps) => {
  const colSpan = () => props.columns.length || 1;

  return (
    <div class="app-table-wrapper">
      <table class="app-table">
        <thead>
          <tr>
            <For each={props.columns}>
              {(column) => (
                <th class={`app-table-header app-table-header--${column.align ?? "left"}`}>
                  {column.header}
                </th>
              )}
            </For>
          </tr>
        </thead>
        <tbody>
          <Show
            when={!props.isLoading && props.data.length > 0}
            fallback={
              <tr>
                <td class="app-table-message" colSpan={colSpan()}>
                  <Show
                    when={props.isLoading}
                    fallback={<span>{props.emptyText ?? "Нет данных"}</span>}
                  >
                    <div class="app-table-loading">
                      <div class="app-table-spinner" />
                      <span>Загрузка...</span>
                    </div>
                  </Show>
                </td>
              </tr>
            }
          >
            <For each={props.data}>
              {(row) => (
                <tr>
                  <For each={props.columns}>
                    {(column) => (
                      <td class={`app-table-cell app-table-cell--${column.align ?? "left"}`}>
                        {column.render ? column.render(row) : (row as any)[column.key]}
                      </td>
                    )}
                  </For>
                </tr>
              )}
            </For>
          </Show>
        </tbody>
      </table>
    </div>
  );
};

