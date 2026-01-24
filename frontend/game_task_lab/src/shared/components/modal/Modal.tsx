import type { JSX } from "solid-js";
import { Show } from "solid-js";
import "./Modal.css";

interface ModalProps {
  isOpen: boolean;
  title: string;
  onClose: () => void;
  children: JSX.Element;
  footer?: JSX.Element;
  showCloseButton?: boolean;
  closeOnOverlayClick?: boolean;
}

export const Modal = (props: ModalProps) => {
  const showCloseButton = props.showCloseButton ?? true;
  const closeOnOverlayClick = props.closeOnOverlayClick ?? true;

  return (
    <Show when={props.isOpen}>
      <div
        class="modal-overlay"
        onClick={() => {
          if (closeOnOverlayClick) props.onClose();
        }}
      >
        <div class="modal" onClick={(e) => e.stopPropagation()}>
          <div class="modal-header">
            <h3>{props.title}</h3>
            {showCloseButton && (
              <button
                class="modal-close"
                onClick={props.onClose}
                aria-label="Закрыть"
              >
                ×
              </button>
            )}
          </div>

          <div class="modal-body">{props.children}</div>

          {props.footer && <div class="modal-footer">{props.footer}</div>}
        </div>
      </div>
    </Show>
  );
};
