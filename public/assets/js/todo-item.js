// TODO refactor to use observedAttributes.

import {html} from '/assets/js/utils.js';

class TodoItem extends HTMLElement {
    constructor() {
        super();
        this.shadow = this.attachShadow({mode: 'open'});
        this.editMode = this.hasAttribute('edit');
        console.log(html)
    }

    async connectedCallback() {
        const id = this.getAttribute('todo-id');
        const datetime = this.getAttribute('todo-datetime');
        const datetimeFormatted = new Intl.DateTimeFormat('et-EE', {
            dateStyle: 'short',
            timeStyle: 'short',
        }).format(new Date(datetime));
        const completed = this.hasAttribute('todo-completed');
        const task = this.getAttribute('todo-task');

        const tpl = `
            <style>
            :host {
              display: block;
              margin: 1rem 0;
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
            <header>
                <time datetime="${datetime}">${datetimeFormatted}</time>
                <div class="actions">
                    <a class="delete" href="/items/${id}">Delete</a>
                    <a class="edit" href="/items/${id}">Edit</a>
                </div>
            </header>
            <form method="PUT" action="/items/${id}" class="task">
                <input type="checkbox" name="Status" ${completed ? 'checked' :''} aria-label="Edit task">
                <span class="task-content">${html(task)}</span>
            </form>
        `;
        this.shadow.innerHTML = tpl;

        this.taskContentElem = this.shadow.querySelector('.task-content');

        this.deleteElem = this.shadow.querySelector('.delete');
        this.deleteElem.addEventListener('click', this.delete.bind(this));

        this.editElem = this.shadow.querySelector('.edit');
        this.editElem.addEventListener('click', this.edit.bind(this));

        this.editFormElem = this.shadow.querySelector('.task');
        this.editFormElem.addEventListener('submit', this.update.bind(this));

        this.statusCheckbox = this.shadow.querySelector('[type="checkbox"]');
        this.statusCheckbox.addEventListener('change', this.update.bind(this));
    }

    async edit(e) {
        e.preventDefault();
        this.editMode = true;
        this.setAttribute('edit', '');
        const id = this.getAttribute('todo-id');
        const listId = this.getAttribute('todo-list-id');

        const editContainer = document.createElement('div');
        editContainer.classList.add('edit-container');

        const input = document.createElement('input');
        input.classList.add('edit-field');
        input.type = 'text';
        input.name = 'Task';
        input.id = `edit-field-${id}`;
        input.required = true;
        input.value = this.taskContentElem.innerText;
        input.focus();
        input.setAttribute('aria-label', 'Edit task');

        const select = document.createElement('select');
        select.classList.add('edit-field', 'select-field');
        select.name = 'ListID';
        select.setAttribute('aria-label', 'Change list');
        select.value = listId;

        (async () => {
            const resp = await fetch('/lists');
            this.lists = await resp.json();
            this.lists.map((l) => {
                const opt = document.createElement('option');
                opt.value = l.ID;
                opt.innerText = l.Name;
                opt.selected = (Number(listId) === l.ID);
                select.appendChild(opt);
            });
        })();

        const hint = document.createElement('label');
        hint.classList.add('edit-hint');
        hint.innerText = 'Press Esc to cancel edit, Enter to save.'
        hint.setAttribute('for', input.id);

        editContainer.append(input);
        editContainer.append(select);
        editContainer.append(hint);

        this.editFormElem.append(editContainer);

        this.taskContentElem.hidden = true;
        e.target.hidden = true;

        this.removeAttribute('todo-completed');

        input.addEventListener('keydown', (e) => {
            if(e.key === 'Escape') {
                this.taskContentElem.hidden = false;
                this.editElem.hidden = false;
                this.editMode = false;
                this.removeAttribute('edit');

                if (this.statusCheckbox.checked) {
                    this.setAttribute('todo-completed', '');
                } else {
                    this.removeAttribute('todo-completed');
                }
                editContainer.remove();
                return;
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

        // If chekckbox is checked then just return.
        if (e.type === 'change') {
            return;
        }

        // If list is change move it to another list.
        const oldListID = this.getAttribute('todo-list-id');
        const newListID = data.get('ListID');
        if (oldListID !== newListID) {
            const newList = document.getElementById(`todo-list-id-${newListID}`);
            this.setAttribute('todo-list-id', newListID);
            newList.addItem(this);
        }

        this.editElem.hidden = false;
        this.taskContentElem.hidden = false;
        if (data.get('Task')) {
            this.taskContentElem.innerText = data.get('Task');
        }
        const container = this.shadow.querySelector('.edit-container');
        container && container.remove();
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
}

customElements.define("todo-item", TodoItem);
