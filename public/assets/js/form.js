const formElem = document.getElementById('items-form');
const itemsElem = document.getElementById('items');
const submitBtn = document.getElementById('submit-btn');

formElem.addEventListener('submit', async (e) => {
    e.preventDefault();
    submitBtn.disabled = true;
    const formData = new FormData(formElem);
    const resp = await fetch(formElem.action, {
        method: formElem.method,
        body: formData,
    });
    const lastID = await resp.text();
    if (resp.ok) {
        const item = document.createElement('todo-item');
        item.setAttribute('todo-id', lastID);
        item.setAttribute('todo-datetime', new Date());
        item.setAttribute('todo-task', formData.get('Task'));
        formData.get('Status') && item.setAttribute('todo-completed', '');
        item.classList.add('-added');
        items.prepend(item);
        formElem.reset();
    } else {
        console.error(resp.statusText);
    }
    submitBtn.disabled = false;
});
