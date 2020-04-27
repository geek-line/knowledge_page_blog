const content = document.getElementById('content')
const input = document.getElementById('input')
const like_button_inline = document.getElementById('like_button_inline')
const like_button_baloon = document.getElementById('like_button_baloon')
const likes_inline = document.getElementById('likes_inline')
const likes_baloon = document.getElementById('likes_baloon')
const knowledge_id = document.getElementById('knowledge_id').value

document.addEventListener('DOMContentLoaded', function () {
   
    if (localStorage.getItem('noLoginLike')) {
        let value = localStorage.getItem('noLoginLike')
        let values = value.split(',')
        for (let i = 0; i < values.length; i++) {
            if (values[i] == knowledge_id) {
                like_button_inline.textContent = 'LIKED'
                like_button_inline.classList.add('liked-button')
                break
            }
        }
    }
    content.innerHTML = input.value
    var p_table_items = document.getElementById("p_table_items"); // 目次を追加する先(table of contents)
    var p_table_items_devise = document.getElementById("p_table_items_devise"); // 目次を追加する先(table of contents)
    var div = document.createElement('div'); // 作成する目次のコンテナ要素
    // .entry-content配下のh2、h3要素を全て取得する
    var matches = document.querySelectorAll('.content h2,.content h3');
    // .entry-content配下のh2、h3要素を全て取得する
    // 取得した見出しタグ要素の数だけ以下の操作を繰り返す
    matches.forEach(function (value, i) {
        // 見出しタグ要素のidを取得し空の場合は内容をidにする
        var id = value.id;
        if (id === '') {
            id = value.textContent;
            value.id = id;
        }
        // 要素がh2タグの場合
        if (value.tagName === 'H2') {
            var ul = document.createElement('ul');
            var li = document.createElement('li');
            var a = document.createElement('a');
            // 追加する<ul><li><a>タイトル</a></li></ul>を準備する
            a.innerHTML = value.textContent;
            a.href = '#' + value.id;
            a.className = "h2 sidenav-close"
            li.appendChild(a)
            ul.appendChild(li);
            // コンテナ要素である<div>の中に要素を追加する
            div.appendChild(ul);
        }
        // 要素がh3タグの場合
        if (value.tagName === 'H3') {
            var ul = document.createElement('ul');
            var li = document.createElement('li');
            var a = document.createElement('a');
            // コンテナ要素である<div>の中から最後の<li>を取得する
            var lastUl = div.lastElementChild;
            var lastLi = lastUl.lastElementChild;
            // 追加する<ul><li><a>タイトル</a></li></ul>を準備する
            a.innerHTML = '&nbsp; ->' + value.textContent;
            a.href = '#' + value.id;
            a.className = "h3 sidenav-close"
            li.appendChild(a)
            ul.appendChild(li);
            // 最後の<li>の中に要素を追加する
            lastLi.appendChild(ul);
        }
    });

    p_table_items_devise.appendChild(div);
    p_table_items.innerHTML = p_table_items_devise.innerHTML
    
    var elems = document.querySelectorAll('.sidenav');
    var instances = M.Sidenav.init(elems,{draggable:true,edge:'right'});
    smoothScroll();
});

const smoothScroll = () =>{
    let links = document.querySelectorAll('.item_devise a[href^="#"]');
    const speed = 200;          // スクロールスピード   
    const divisor = 100;        // 分割数
    const tolerance = 5;        // 許容誤差
    const headerHeight = 40;     // 固定ヘッダーがある場合はその高さ分ずらす
    const interval = speed/divisor;
    for(let i = 0; i < links.length; i++){
      links[i].addEventListener('click',(e)=>{
        e.preventDefault();
        let nowY = window.pageYOffset;
        const href = e.currentTarget.getAttribute('href');   //href取得
        const splitHref = href.split('#');
        const targetID = splitHref[1];
        const target = document.getElementById(targetID);  
        if( target != null){
          const targetRectTop = target.getBoundingClientRect().top;
          const targetY = targetRectTop + nowY - headerHeight;
          const minY = Math.abs((targetY - nowY)/divisor);
          doScroll(minY,nowY,targetY,tolerance,interval);
        }
      });
    }
}
  
const doScroll = (minY,nowY,targetY,tolerance,interval) =>{
    let toY ;
    if( targetY < nowY ){
        toY = nowY - minY;
    }else{
        toY = nowY + minY;
    }
    window.scrollTo(0, toY);
    if( targetY - tolerance > toY || toY > targetY + tolerance){
    window.setTimeout(doScroll,interval,minY,toY,targetY,tolerance,interval);
    }else{
    return false;
    }
}

function sendLikeFromBaloon() {
    let values = []
    let value = ''
    if (localStorage.getItem('noLoginLike')) {
        value = localStorage.getItem('noLoginLike')
        values = value.split(',')
    }
    let isFound = false
    for (let i = 0; i < values.length; i++) {
        if (values[i] == knowledge_id) {
            values.splice(i, 1)
            value = values.join()
            isFound = true
            break
        }
    }
    if (!isFound) {
        values.push(knowledge_id)
        value = values.join()
    }
    const XHR = new XMLHttpRequest()
    let formdata = new FormData(document.getElementById('like_form_baloon'))
    if (isFound) {
        XHR.open('PUT', '/knowledges/like')
    } else {
        XHR.open('POST', '/knowledges/like')
    }
    XHR.onreadystatechange = function () {
        if (XHR.readyState === 4) {
            if (XHR.status === 200) {
                if (isFound) {
                    likes_baloon.textContent = Number(likes_baloon.textContent) - 1
                    like_button_baloon.textContent = 'LIKE'
                    like_button_baloon.classList.remove('liked-button')
                } else {
                    likes_baloon.textContent = Number(likes_baloon.textContent) + 1
                    like_button_baloon.textContent = 'LIKED'
                    like_button_baloon.classList.add('liked-button')
                }
                localStorage.setItem('noLoginLike', value)
            } else {
                alert('データが正常に送れませんでした')
            }
        }
    }
    XHR.send(formdata)
}
function sendLikeFromInline() {
    let values = []
    let value = ''
    if (localStorage.getItem('noLoginLike')) {
        value = localStorage.getItem('noLoginLike')
        values = value.split(',')
    }
    let isFound = false
    for (let i = 0; i < values.length; i++) {
        if (values[i] == knowledge_id) {
            values.splice(i, 1)
            value = values.join()
            isFound = true
            break
        }
    }
    if (!isFound) {
        values.push(knowledge_id)
        value = values.join()
    }
    const XHR = new XMLHttpRequest()
    let formdata = new FormData(document.getElementById('like_form_inline'))
    if (isFound) {
        XHR.open('PUT', '/knowledges/like')
    } else {
        XHR.open('POST', '/knowledges/like')
    }
    XHR.onreadystatechange = function () {
        if (XHR.readyState === 4) {
            if (XHR.status === 200) {
                if (isFound) {
                    likes_inline.textContent = Number(likes_inline.textContent) - 1
                    like_button_inline.textContent = 'LIKE'
                    like_button_inline.classList.remove('liked-button')
                } else {
                    likes_inline.textContent = Number(likes_inline.textContent) + 1
                    like_button_inline.textContent = 'LIKED'
                    like_button_inline.classList.add('liked-button')
                }
                localStorage.setItem('noLoginLike', value)
            } else {
                alert('データが正常に送れませんでした')
            }
        }
    }
    XHR.send(formdata)
}