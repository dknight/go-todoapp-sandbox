class TodoListLink extends HTMLElement {
    constructor() {
        super();
        this.shadow = this.attachShadow({mode: 'open'});
    }

    get styles() {
        return `
            <style>
            :host {
              display: block;
              padding: .1rem .2rem;
            }
            :host(.-deleted) {
              background-color: #c00;
              transition: all .8s ease-out;
              opacity: 0;
              border-color: #f00;
            }
            .icon {
              all: unset;
              width: 1.4rem;
              height: 1.4rem;
              transform: translateY(-25%);
              cursor: pointer;
            }
            .edit {
              background: transparent url('/assets/img/edit.svg') center;
              background-size: 100%;
              filter: brightness(.9)
                invert(.7)
                sepia(.5)
                hue-rotate(100deg)
                saturate(200%);
            }
            .delete {
              background: transparent url('/assets/img/delete.svg') center;
              background-size: 100%;
              filter: brightness(.9)
                invert(.7)
                sepia(.5)
                hue-rotate(330deg)
                saturate(3000%);
            }
            .edit-field {
              font-size: inherit;
              font-family: inherit;
              padding: .5rem;
              color: var(--text);
              background-color: var(--bg);
              border: 1px solid var(--border);
              border-radius: 5px;
              box-shadow: none;
              max-width: 100%;
              width: 100%;
              display: inline-block;
            }
            .edit-hint {
              color: var(--text-light);
              display: block;
              font-size: .9rem;
              margin-top: .2rem;
            }
            </style>`;
    }

    get layouts() {
        return {
            default: `
                <a href="/#todo-list-id-${this.listID}">${this.listName} (${this.totalCount})</a>
                <button class="icon edit" title="Edit"></button>
                <button class="icon delete" title="Delete"></button>`,
            edit: `
                <form method="PUT" action="/lists/${this.listID}">
                    <input type="text"
                           class="edit-field"
                           value="${this.listName}"
                           name="Name"
                           id="List-Name-${this.listID}">
                    <label for="${this.listID}" class="edit-hint">
                      Press Esc to cancel edit, Enter to save.
                    </label>
                </form>`,
        };
    }

    static get observedAttributes() {
        return ['edit'];
    }

    render() {
        this.listID = this.getAttribute('todo-list-id');
        this.listName = this.getAttribute('todo-list-name');
        this.editMode = this.hasAttribute('edit');
        this.totalCount = Number(this.getAttribute('todo-list-count'));
        this.completedCount = Number(this.getAttribute('todo-list-completed'));

        this.shadow.innerHTML = this.styles
            + (this.editMode ? this.layouts.edit : this.layouts.default);
    }

    hydrate() {
        if (this.editMode) {
            const form = this.shadow.querySelector('form');
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                const data = new FormData(form);
                const resp = await fetch(form.action, {
                    method: 'PUT',
                    body: data,
                });
                if (resp.ok) {
                    this.setAttribute('todo-list-name', data.get('Name'));
                    this.removeAttribute('edit');
                } else {
                    console.error(resp.statusText);
                }
            });

            form.addEventListener('keydown', (e) => {
                if(e.key === 'Escape') {
                    this.removeAttribute('edit');
                }
            });
        } else {
            const editBtn = this.shadow.querySelector('button.edit');
            editBtn.addEventListener('click', (e) => {
                this.setAttribute('edit', '');
                const input = this.shadow.getElementById(`List-Name-${this.listID}`);
                input && input.focus();
            });
            const deleteBtn = this.shadow.querySelector('button.delete');
            deleteBtn.addEventListener('click', async (e) => {
                const ok = window.confirm(
                    'Are you sure? All items will be removed!'
                );
                if (!ok) {
                    return;
                }
                const resp = await fetch(`/lists/${this.listID}`, {
                    method: 'DELETE',
                });
                if (resp.ok) {
                    this.classList.add('-deleted');
                    this.addEventListener('transitionend', (e) => {
                        this.parentElement.remove();
                    });
                } else {
                    console.error(resp.statusText);
                }
            });
        }
    }

    attributeChangedCallback(oldValue, newValue) {
        if (oldValue === newValue) {
            return;
        }
        this.render();
        this.hydrate();
    }

    connectedCallback() {
        this.render();
        this.hydrate();
    }
}

customElements.define("todo-list-link", TodoListLink);
