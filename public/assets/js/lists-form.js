const listFormElem = document.getElementById('lists-form');
const listsElem = document.getElementById('lists');
const submitBtn = document.getElementById('submit-btn');

listFormElem.addEventListener('submit', async (e) => {
    e.preventDefault();
    submitBtn.disabled = true;
    const formData = new FormData(listFormElem);
    const resp = await fetch('/lists', {
        method: 'POST',
        body: formData,
    });
    const lastID = await resp.text();
    if (resp.ok) {
        const link = document.createElement('todo-list-link');
        link.setAttribute('todo-list-id', lastID);
        link.setAttribute('todo-list-name', formData.get('Name'));

        const li = document.createElement('li');
        li.appendChild(link);
        listsElem.append(li);

        listFormElem.reset();
    } else {
        console.error(resp.statusText)
    }
    submitBtn.disabled = false;
});
