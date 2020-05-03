const search_input = document.getElementById('search_input')
const search_submit = document.getElementById('search_submit')


function go(){
    //EnterキーならSubmit
    if(window.event.keyCode==13)document.getElementById("search_submit").click();
}

search_submit.addEventListener('click', function (e) {
    if (!search_input.value) {
        e.preventDefault()
        return
    }
    const XHR = new XMLHttpRequest()
    XHR.open('GET', '/search?q=' + search_input.value.replace(/\s+/g, '&'))
    XHR.onreadystatechange = function () {
        if (XHR.readyState === 4) {
            if (XHR.status === 200) {
                location.href = '/search?q=' + search_input.value.replace(/\s+/g, '&')
            } else {
                alert('キーワードを正常に送信できませんでした。')
            }
        }
    }
    XHR.send()
})