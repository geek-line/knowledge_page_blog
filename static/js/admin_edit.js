const form = document.getElementById('form')
const select_eyecatch = document.getElementById('select_eyecatch')
const file_preview = document.getElementById('file_preview')
const add_tag_button = document.getElementById('add_tag_button')
window.addEventListener('DOMContentLoaded', function () {
    select_display.textContent = null;
    let selectedTags = document.getElementsByName('selectedTagsID')[0].value
    selectedTags = selectedTags.split(',')
    selectedTags.pop()
    for (let i = 0; i < selectedTags.length; i++) {
        const selectElement = document.createElement('div')
        selectElement.innerHTML = add_tag_button.value
        document.getElementById('select_display').appendChild(selectElement)
        selectElement.childNodes[1].addEventListener('click', function () {
            selectElement.parentNode.removeChild(selectElement);
        })
        for (let j = 0; j < selectElement.childNodes[0].length; j++) {
            if (selectElement.childNodes[0][j].value === selectedTags[i]) {
                selectElement.childNodes[0][j].selected = true
                break
            }
        }
    }
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
function sendData() {
    const content = document.getElementById('tinymce_body_ifr').contentWindow.document.getElementById('tinymce').innerHTML
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
    formdata.append('tags', tags)
    const XHR = new XMLHttpRequest()
    XHR.open('PUT', '/admin/save/')
    XHR.send(formdata)
    XHR.onreadystatechange = function () {
        if (XHR.readyState === 4) {
            if (XHR.status === 200) {
                alert('データが更新されました')
                location.href = "/admin/knowledges";
            } else {
                alert('データが正常に送れませんでした')
            }
        }
    }
}