// TODO refactor to use observedAttributes.

import {html} from '/assets/js/utils.js';

class TodoItem extends HTMLElement {
    constructor() {
        super();
        this.shadow = this.attachShadow({mode: 'open'});
    }

    static get observedAttributes() {
        return ['edit'];
    }

    get styles() {
        return `
            <style>
            :host {
              display: block;
              padding: 1rem;
              border-radius: 5px;
            }
            :host(.-added) {
              animation: fade-in .8s ease-out;
            }
            :host(.-deleted) {
              background-color: #c00;
              transition: all .8s ease-out;
              opacity: 0;
              border-color: #f00;
            }
            :host([todo-completed]:not([edit])) {
              opacity: .5;
            }
            :host([todo-completed]) .task-content {
              text-decoration: line-through;
            }
            .actions {
              display: flex;
              gap: .5rem;
              visibility: hidden;
            }
            :host(:hover) .actions {
                visibility: visible;
            }
            header {
              display: flex;
              gap: .7rem;
            }
            time {
              color: var(--text-light);
            }
            .task {
              display: flex;
              align-items: flex-start;
              gap: .7rem;
            }
            input[type="checkbox"] {
              min-width: 1.5rem;
              min-height: 1.5rem;
            }
            .task-content {
              margin-top: .15rem;
            }
            .icon {
              width: 1.4rem;
              height: 1.4rem;
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
            .edit-container {
              flex: 1;
              margin-right: 5%;
            }
            .edit-hint {
              color: var(--text-light);
              display: block;
              font-size: .9rem;
              margin-top: .2rem;
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
            .select-field {
              margin-top: .2rem;
              width: auto;
            }
            .-error {
              color: #f30;
            }
            </style>
        `;
    }

    get layouts() {
        return {
            default: `
            <header>
                <time datetime="${this.datetime}">
                    ${this.datetimeFormatted}
                </time>
                <div class="actions">
                    <a class="icon delete" href="/items/${this.todoID}"></a>
                    <a class="icon edit" href="/items/${this.todoID}"></a>
                </div>
            </header>
            <form method="PUT" action="/items/${this.todoID}" class="task">
                <input type="checkbox" name="Status"
                    ${this.completed ? 'checked' :''} aria-label="Edit task">
                <span class="task-content">${html(this.task)}</span>
            </form>`,
            edit: `
                <header>
                    <time datetime="${this.datetime}">
                        ${this.datetimeFormatted}
                    </time>
                    <div class="actions">
                        <a class="delete" href="/items/${this.todoID}">Delete</a>
                    </div>
                </header>
                <form method="PUT" action="/items/${this.todoID}" class="task">
                    <input type="checkbox"
                           name="Status"
                           aria-label="Edit task"
                           ${this.completed ? 'checked' :''}>
                    <div class="edit-container">
                        <input type="text"
                            name="Task"
                            id="edit-field-${this.todoID}"
                            value="${this.task}"
                            aria-label="Task name"
                            class="edit-field"
                            required>
                        <select class="edit-field select-field"
                            name="ListID" aria-label="Change list">
                        </select>
                        <label class="edit-hint" for="edit-field-${this.todoID}">
                            Press Esc to cancel edit, Enter to save.
                        </label>
                    </div>
                </form>
            `
        };
    }

    render() {
        this.editMode = this.hasAttribute('edit');
        this.todoID = this.getAttribute('todo-id');
        this.listID = Number(this.getAttribute('todo-list-id'));
        this.datetime = this.getAttribute('todo-datetime');
        this.datetimeFormatted = new Intl.DateTimeFormat('et-EE', {
            dateStyle: 'short',
            timeStyle: 'short',
        }).format(new Date(this.datetime));
        this.completed = this.hasAttribute('todo-completed');
        this.task = this.getAttribute('todo-task');

        this.shadow.innerHTML = this.styles +
            (this.editMode ? this.layouts.edit : this.layouts.default);
    }

    hydrate() {
        this.deleteElem = this.shadow.querySelector('.delete');
        this.deleteElem.addEventListener('click', this.delete.bind(this));

        this.editElem = this.shadow.querySelector('.edit');
        if (this.editElem) {
            this.editElem.addEventListener('click', this.edit.bind(this));
        }

        this.editFormElem = this.shadow.querySelector('.task');
        if (this.editFormElem) {
            this.editFormElem.addEventListener('submit', this.update.bind(this));
        }

        this.statusCheckbox = this.shadow.querySelector('[type="checkbox"]');
        this.statusCheckbox.addEventListener('change', this.update.bind(this));
    }

    async edit(e) {
        e.preventDefault();
        this.setAttribute('edit', '');

        const select = this.shadow.querySelector('select');
        select.value = this.listID;

        const input = this.shadow.querySelector('[name="Task"]');

        (async () => {
            if (!this.lists) {
                const resp = await fetch('/lists');
                this.lists = await resp.json();
            }
            this.lists.map((l) => {
                const opt = document.createElement('option');
                opt.value = l.ID;
                opt.innerText = l.Name;
                opt.selected = (this.listID === l.ID);
                select.appendChild(opt);
            });
        })();

        input.addEventListener('keydown', (e) => {
            if(e.key === 'Escape') {
                this.removeAttribute('edit');

                this.statusCheckbox.checked
                    ? this.setAttribute('todo-completed', '')
                    : this.removeAttribute('todo-completed');
            }
        });
    }

    async update(e) {
        e.preventDefault();
        const data = new FormData(this.editFormElem);
        const resp = await fetch(this.editFormElem.action, {
            method: 'PUT',
            body: data,
        });
        if (!resp.ok) {
            const errMsg = 'Error occured. Try again later.';
            let errElem = document.querySelector('-error');
            if (!errElem) {
                errElem = document.createElement('span');
                errElem.classList.add('-error');
                errElem.innerText = errMsg;
                this.editFormElem.insertAdjacentElement('beforebegin', errElem);
            }
            console.error(errMsg);
        }

        this.statusCheckbox.checked
            ? this.setAttribute('todo-completed', '')
            : this.removeAttribute('todo-completed');

        // If checkbox is checked then just return.
        if (e.type === 'change') {
            return;
        }

        if (data.get('Task')) {
            this.setAttribute('todo-task', data.get('Task'));
        }

        // If list is change move it to another list.
        const oldListID = this.getAttribute('todo-list-id');
        const newListID = data.get('ListID');
        if (oldListID !== newListID) {
            const newList = document.getElementById(`todo-list-id-${newListID}`);
            this.setAttribute('todo-list-id', newListID);
            newList.addItem(this);
        }

        this.removeAttribute('edit');
    }

    async delete(e) {
        e.preventDefault();
        if (!confirm("Are you sure?")) {
            return;
        }
        const target = e.target;
        const href = target.href;
        const resp = await fetch(href, {
            method: 'DELETE',
        });
        if (resp.ok) {
            this.classList.add('-deleted');
            this.addEventListener('transitionend', (evt) => evt.target.remove());
        } else {
            console.error(resp.statusText);
        }
    }

    connectedCallback() {
        this.render();
        this.hydrate();
    }

    attributeChangedCallback(oldValue, newValue) {
        if (oldValue === newValue) {
            return;
        }
        this.render();
        this.hydrate();
    }
}

customElements.define("todo-item", TodoItem);
