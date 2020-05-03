const search_input = document.getElementById('search_input')
const search_submit = document.getElementById('search_submit')

search_submit.addEventListener('click', function (e) {
    if (!search_input.value) {
        e.preventDefault()
        return
    }
    const XHR = new XMLHttpRequest()
    let queries = search_input.value.split(/\s+/g)
    for (let i = 0; i < queries.length; i++){
        queries[i] = encodeURIComponent(queries[i])
    }
    const qvalue = queries.join('+')
    console.log(qvalue)
    XHR.open('GET', '/search?q=' + qvalue)
    XHR.onreadystatechange = function () {
        if (XHR.readyState === 4) {
            if (XHR.status === 200) {
                location.href = '/search?q=' + qvalue
            } else {
                alert('キーワードを正常に送信できませんでした。')
            }
        }
    }
    XHR.send()
})