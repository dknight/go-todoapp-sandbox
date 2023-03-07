class TodoList extends HTMLElement {
    constructor() {
        super();
        this.shadow = this.attachShadow({mode: 'open'});
    }

    async render() {
        this.listID = this.getAttribute('todo-list-id');
        this.listName = this.getAttribute('todo-list-name');
        const resp = await fetch(`/items/${this.listID}`);
        const items = await resp.json();

        this.setAttribute('role', 'list');
        const styles = `
        <style>
          :host {

          }
        </style>`;

        let tpl = `<h3 class="heading">${this.listName}</h3>`;
        if (items && items.length) {
            tpl += items.reduce((s, item) => s + this._renderItem(item), '');
        }
        tpl += '<hr>';
        this.shadow.innerHTML = styles + tpl;
    }

    _renderItem(item) {
        return `<todo-item
            role="listitem"
            todo-id="${item.ID}"
            todo-list-id="${item.ListID}"
            todo-datetime="${item.CreatedAt}"
            todo-task="${item.Task}"
            ${item.Status ? 'todo-completed' : ''}
        ></todo-item>`;
    }

    addItem(item) {
        const heading = this.shadow.querySelector('.heading');
        heading.insertAdjacentElement('afterend', item);
    }


    hydrate() {}

    connectedCallback() {
        this.render();
        this.hydrate();
    }
}

customElements.define('todo-list', TodoList);
