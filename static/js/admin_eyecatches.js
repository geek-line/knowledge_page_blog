const submit_btn_post = document.getElementById('submit_btn_post')
const submit_btn_put = document.querySelectorAll('#submit_btn_put')
const submit_btn_delete = document.querySelectorAll('#submit_btn_delete')
const file_uploader_post = document.getElementById('file_uploader_post')
const file_uploader_put = document.querySelectorAll('#file_uploader_put')
const file_previews = document.querySelectorAll('#file_preview')
const forms = document.querySelectorAll('#form')
submit_btn_post.addEventListener('click', function () {
    const timestamp = new Date().getTime();
    const filename = "file-" + timestamp + '-' + file_uploader_post.files[0].name;
    s3.putObject({ Key: 'eyecatches/' + filename, ContentType: file_uploader_post.files[0].type, Body: file_uploader_post.files[0], ACL: "public-read" },
        function (err, data) {
            if (data !== null) {
                document.getElementById('src_post').value = "https://knowledge-blog.s3-ap-northeast-1.amazonaws.com/" + 'eyecatches/' + filename
                let formdata = new FormData(document.getElementById('form_post'))
                const XHR = new XMLHttpRequest()
                XHR.open('POST', '/admin/eyecatches/')
                XHR.send(formdata)
                XHR.onreadystatechange = function () {
                    if (XHR.readyState === 4) {
                        if (XHR.status === 200) {
                            alert('データが更新されました')
                            location.href = "/admin/eyecatches/";
                        } else {
                            alert('データが正常に送れませんでした')
                        }
                    }
                }
            } else {
                alert("アップロード失敗.");
            }
        }
    )
})
for (let i = 0; i < file_uploader_put.length; i++) {
    file_uploader_put[i].addEventListener('change', function (e) {
        const file = e.target.files[0];
        // ファイルのブラウザ上でのURLを取得する
        const blobUrl = window.URL.createObjectURL(file);
        file_previews[i].src = blobUrl
    })
}
for (let i = 0; i < forms.length; i++) {
    submit_btn_put[i].addEventListener('click', function () {
        if (file_uploader_put[i].files[0] != null) {
            const timestamp = new Date().getTime();
            const filename = "file-" + timestamp + '-' + file_uploader_put[i].files[0].name;
            let key = document.querySelectorAll('#src_put')[i].value.replace(/^https:\/\/knowledge-blog\.s3-ap-northeast-1\.amazonaws\.com\//, '')
            s3.putObject({ Key: 'eyecatches/' + filename, ContentType: file_uploader_put[i].files[0].type, Body: file_uploader_put[i].files[0], ACL: "public-read" },
                function (err, data) {
                    if (data !== null) {
                        document.querySelectorAll('#src_put')[i].value = "https://knowledge-blog.s3-ap-northeast-1.amazonaws.com/" + 'eyecatches/' + filename
                        let formdata = new FormData(forms[i])
                        const XHR = new XMLHttpRequest()
                        XHR.open('PUT', '/admin/eyecatches/')
                        XHR.send(formdata)
                        XHR.onreadystatechange = function () {
                            if (XHR.readyState === 4) {
                                if (XHR.status === 200) {
                                    s3.deleteObject({ Key: key }, function (err, data) {
                                        if (err != null) {
                                            alert('データの削除に失敗しました')
                                            return
                                        }
                                        alert('データが更新されました')
                                        location.href = "/admin/eyecatches/";
                                    })
                                } else {
                                    alert('データが正常に送れませんでした')
                                    return
                                }
                            }
                        }
                    } else {
                        alert("アップロード失敗.");
                    }
                }
            )
        } else {
            const formdata = new FormData(forms[i])
            const XHR = new XMLHttpRequest()
            XHR.open('PUT', '/admin/eyecatches/')
            XHR.send(formdata)
            XHR.onreadystatechange = function () {
                if (XHR.readyState === 4) {
                    if (XHR.status === 200) {
                        alert('データが更新されました')
                        location.href = "/admin/eyecatches/";
                    } else {
                        alert('データが正常に送れませんでした')
                    }
                }
            }
        }
    })
    submit_btn_delete[i].addEventListener('click', function () {
        const formdata = new FormData(forms[i])
        const XHR = new XMLHttpRequest()
        let key = document.querySelectorAll('#src_put')[i].value.replace(/^https:\/\/knowledge-blog\.s3-ap-northeast-1\.amazonaws\.com\//, '')
        s3.deleteObject({ Key: key }, function (err, data) {
            XHR.open('DELETE', '/admin/eyecatches/')
            XHR.send(formdata)
            XHR.onreadystatechange = function () {
                if (XHR.readyState === 4) {
                    if (XHR.status === 200) {
                        alert('データが更新されました')
                        location.href = "/admin/eyecatches/";
                    } else {
                        alert('データが正常に送れませんでした')
                    }
                }
            }
        })
    })
}