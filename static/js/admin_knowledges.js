const delete_buttons = document.querySelectorAll('#delete_button')

for (let i = 0; i < delete_buttons.length; i++){
    delete_buttons[i].addEventListener('click', function () {
        const result = window.confirm('本当に削除してもよろしいですか？(削除すると2度と編集できなくなります)')
        if (result) {
            const XHR = new XMLHttpRequest()
            XHR.open('DELETE', '/admin/delete/' + delete_button[i].value)
            XHR.onreadystatechange = function () {
                if (XHR.readyState === 4) {
                    if (XHR.status === 200) {
                        location.href = "/admin/knowledges/";
                    } else {
                        alert('データが正常に送れませんでした')
                    }
                }
            }
            XHR.send()
        } else {
            return
        }
    })
}

