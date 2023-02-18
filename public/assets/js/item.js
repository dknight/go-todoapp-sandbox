class TodoItem extends HTMLElement {
    constructor() {
        super();
        this.shadow = this.attachShadow({mode: 'open'});
    }

    connectedCallback() {
        const id = this.getAttribute('todo-id');
        const datetime = this.getAttribute('todo-datetime');
        const datetimeFormatted = new Intl.DateTimeFormat('et-EE', {
            dateStyle: 'short',
            timeStyle: 'medium',
        }).format(new Date(datetime));
        const completed = this.hasAttribute('todo-completed');
        const task = this.getAttribute('todo-task');

        const tpl = `
            <style>
            :host {
              display: block;
              border: 1px solid var(--border);
              padding: 1rem;
              margin-bottom: 1rem;
              border-radius: 5px;
            }
            :host(:not([todo-completed]):hover) {
              border-color: var(--code);
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
            :host([todo-completed]) {
              opacity: .5;
            }
            :host([todo-completed]) .task-content {
              text-decoration: line-through;
            }
            time {
              color: var(--text-light);
            }
            .task {
              display: flex;
              align-items: flex-start;
              gap: .7rem;
              margin-block-start: .5rem;
              margin-block-end: .5rem;
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
            .-error {
              color: #f30;
            }
            footer {
              display: flex;
              gap: 1rem;
            }
            </style>
            <header>
                <time datetime="${datetime}">${datetimeFormatted}</time>
            </header>
            <form method="PUT" action="/items/${id}" class="task">
                <input type="checkbox" name="Status" ${completed ? 'checked' :''}>
                <span class="task-content">${task}</span>
            </form>
            <footer>
                <a class="delete" href="/items/${id}">Delete</a>
                <a class="edit" href="/items/${id}">Edit</a>
            </footer>
        `;
        this.shadow.innerHTML = tpl;
        
        this.taskContentElem = this.shadow.querySelector('.task-content');
        
        this.deleteElem = this.shadow.querySelector('.delete');
        this.deleteElem.addEventListener('click', this.delete.bind(this));

        this.editElem = this.shadow.querySelector('.edit');
        this.editElem.addEventListener('click', this.createEditField.bind(this));

        this.editFormElem = this.shadow.querySelector('.task');
        this.editFormElem.addEventListener('submit', this.update.bind(this));

        this.statusCheckbox = this.shadow.querySelector('[type="checkbox"]');
        this.statusCheckbox.addEventListener('change', (e) => {
            this.update(e);
            if (e.target.checked) {
                this.setAttribute('todo-completed', '');
            } else {
                this.removeAttribute('todo-completed');
            }
        });
    }

    createEditField(e) {
        e.preventDefault();
        const id = this.getAttribute('todo-id');

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

        const hint = document.createElement('label');
        hint.classList.add('edit-hint');
        hint.innerText = 'Press Esc to cancel edit, Enter to save.'
        hint.setAttribute('for', input.id);

        editContainer.append(input);
        editContainer.append(hint);

        this.editFormElem.append(editContainer);

        this.taskContentElem.hidden = true;
        e.target.hidden = true;

        input.addEventListener('keydown', (e) => {
            if(e.key === 'Escape') {
                this.taskContentElem.hidden = false;
                this.editElem.hidden = false;
                editContainer.remove();
                return;
            }
        });
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
        this.editElem.hidden = false;
        this.taskContentElem.hidden = false;
        if (data.get('Task')) {
            this.taskContentElem.innerText = data.get('Task');
        }
        const container = this.shadow.querySelector('.edit-container');
        container && container.remove();
    }
}

customElements.define("todo-item", TodoItem);
