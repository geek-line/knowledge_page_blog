tinymce.init({
    selector: '#tinymce_body',
    branding: false, // クレジットの削除
    height: "640",
    tinycomments_mode: 'embedded',
    tinycomments_author: 'Author name',
    plugins: 'link image lists table codesample',
    toolbar: 'undo redo | styleselect | link bold italic | image codesample | numlist bullist | table tabledelete',
    images_upload_handler: function (blobInfo, success, failure) {
        setTimeout(function () {
            const file = blobInfo.blob()
            const timestamp = new Date().getTime();
            const filename = "file" + timestamp + file.name;
            s3.putObject({ Key: 'uploads/' + filename, ContentType: blobInfo.blob().type, Body: blobInfo.blob(), ACL: "public-read" },
                function (err, data) {
                    if (data !== null) {
                        const srcHTML = "https://knowledge-blog.s3-ap-northeast-1.amazonaws.com/" + 'uploads/' + filename
                        success(srcHTML);
                    } else {
                        alert("アップロード失敗.");
                    }
                }
            );
        }, 2000);
    }
});