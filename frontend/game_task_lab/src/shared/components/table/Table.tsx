import type { JSX } from "solid-js";
import { createMemo, For, Show } from "solid-js";
import {
  createSolidTable,
  getCoreRowModel,
  flexRender,
  type ColumnDef,
  type Table as TanStackTable,
} from "@tanstack/solid-table";
import "./Table.css";

export interface TableColumn<T = any> {
  key: string;
  header: string;
  align?: "left" | "center" | "right";
  render?: (row: T) => JSX.Element;
  accessorKey?: string;
  accessorFn?: (row: T) => any;
}

interface TableProps<T = any> {
  columns: TableColumn<T>[];
  data: T[];
  isLoading?: boolean;
  emptyText?: string;
  selectedRowId?: string | null;
  onRowClick?: (row: T) => void;
  getRowId?: (row: T) => string;
}

export const Table = <T extends Record<string, any>>(props: TableProps<T>) => {
  const columnDefs = createMemo<ColumnDef<T>[]>(() =>
    props.columns.map((col) => ({
      id: col.key,
      header: col.header,
      accessorKey: col.accessorKey || col.key,
      accessorFn: col.accessorFn,
      cell: (info: any) => {
        const row = info.row.original;
        if (col.render) {
          return col.render(row);
        }
        return info.getValue();
      },
    }))
  );

  const table = createMemo(() =>
    createSolidTable({
      get data() {
        return props.data;
      },
      columns: columnDefs(),
      getCoreRowModel: getCoreRowModel(),
      getRowId: props.getRowId || ((row: T) => (row as any).id?.toString() || ""),
      enableRowSelection: false,
    })
  );

  const colSpan = () => props.columns.length || 1;

  return (
    <div class="app-table-wrapper">
      <table class="app-table">
        <thead>
          <For each={table().getHeaderGroups()}>
            {(headerGroup) => (
              <tr>
                <For each={headerGroup.headers}>
                  {(header) => {
                    const column = props.columns.find((col) => col.key === header.id);
                    const align = column?.align ?? "left";
                    return (
                      <th
                        class={`app-table-header app-table-header--${align}`}
                        style={{ width: header.getSize() !== 150 ? `${header.getSize()}px` : undefined }}
                      >
                        {header.isPlaceholder
                          ? null
                          : flexRender(header.column.columnDef.header, header.getContext())}
                      </th>
                    );
                  }}
                </For>
              </tr>
            )}
          </For>
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
            <For each={table().getRowModel().rows}>
              {(row) => {
                const rowData = row.original;
                const rowId = props.getRowId
                  ? props.getRowId(rowData)
                  : (rowData as any).id?.toString() || "";
                const isSelected = () => {
                  const selectedId = props.selectedRowId?.toString() || null;
                  return selectedId === rowId;
                };

                return (
                  <tr
                    class={isSelected() ? "app-table-row--selected" : ""}
                    onClick={() => props.onRowClick?.(rowData)}
                    style={{ cursor: props.onRowClick ? "pointer" : "default" }}
                  >
                    <For each={row.getVisibleCells()}>
                      {(cell) => {
                        const column = props.columns.find((col) => col.key === cell.column.id);
                        const align = column?.align ?? "left";
                        return (
                          <td class={`app-table-cell app-table-cell--${align}`}>
                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                          </td>
                        );
                      }}
                    </For>
                  </tr>
                );
              }}
            </For>
          </Show>
        </tbody>
      </table>
    </div>
  );
};

