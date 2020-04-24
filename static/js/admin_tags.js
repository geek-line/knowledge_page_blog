const forms = document.querySelectorAll('#form')
const update_buttons = document.querySelectorAll('#update_button')
const delete_buttons = document.querySelectorAll('#delete_button')
for (let i = 0; i < form.length; i++) {
    update_buttons[i].addEventListener('click', function () {
        const formdata = new FormData(forms[i])
        const XHR = new XMLHttpRequest()
        XHR.open('PUT', '/admin/tags/')
        XHR.send(formdata)
        XHR.onreadystatechange = function () {
            if (XHR.readyState === 4) {
                if (XHR.status === 200) {
                    alert('データが更新されました')
                    location.href = "/admin/tags/";
                } else {
                    alert('データが正常に送れませんでした')
                }
            }
        }
    })
    delete_buttons[i].addEventListener('click', function () {
        const formdata = new FormData(forms[i])
        const XHR = new XMLHttpRequest()
        XHR.open('DELETE', '/admin/tags/')
        XHR.send(formdata)
        XHR.onreadystatechange = function () {
            if (XHR.readyState === 4) {
                if (XHR.status === 200) {
                    alert('データが更新されました')
                    location.href = "/admin/tags/";
                } else {
                    alert('データが正常に送れませんでした')
                }
            }
        }
    })
}