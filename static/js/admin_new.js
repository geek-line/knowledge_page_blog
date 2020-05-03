const form = document.getElementById('form')
const add_tag_button = document.getElementById('add_tag_button')
const submit_button = document.getElementById('submit_button')
const select_eyecatch = document.getElementById('select_eyecatch')
const file_preview = document.getElementById('file_preview')
window.addEventListener('DOMContentLoaded', function () {
    select_display.textContent = null;
    file_preview.src = select_eyecatch.value
})
add_tag_button.addEventListener('click', function () {
    const selectElement = document.createElement('div')
    selectElement.innerHTML = add_tag_button.value
    document.getElementById('select_display').appendChild(selectElement)
    selectElement.childNodes[1].addEventListener('click', function () {
        selectElement.parentNode.removeChild(selectElement);
    })
})
select_eyecatch.addEventListener('change', function (e) {
    file_preview.src = e.target.value
})
submit_button.addEventListener('click', function (e) {
    const content = document.getElementById('tinymce_body_ifr').contentWindow.document.getElementById('tinymce').innerHTML
    const rowContent = content.replace(/<("[^"]*"|'[^']*'|[^'">])*>/g, '').replace(/\n/g, '')
    const elem_tags = document.getElementsByClassName('elem_tag')
    if (document.getElementById('form-title').value == '') {
        alert('タイトルを入力してください')
        e.preventDefault()
        return
    }
    let arr = {}
    let tags = ''
    for (let i = 0; i < elem_tags.length; i++) {
        if (arr[elem_tags[i].value]) {
            alert('タグが重複しています')
            e.preventDefault()
            return
        }
        arr[elem_tags[i].value] = true
        tags += elem_tags[i].value
        tags += ','
    }
    tags = tags.slice(0, -1)
    let formdata = new FormData(form)
    formdata.append('content', content)
    formdata.append('row_content', rowContent)
    formdata.append('tags', tags)
    const XHR = new XMLHttpRequest()
    XHR.open('POST', '/admin/save/')
    XHR.send(formdata)
    XHR.onreadystatechange = function () {
        if (XHR.readyState === 4) {
            if (XHR.status === 200) {
                location.href = "/admin/knowledges";
            } else {
                alert('データが正常に送れませんでした')
            }
        }
    }
});