# knowledge_page_blog

## デプロイ先
https://code-database.com/

## 作業環境 
Go(1.14 darwin/amd64)  
MySQL(8.0.19)

## 使用ライブラリ
 - joho/godotenv
 - gorilla/sessions
 - go-sql-driver/mysql

## 作業ルール  
 - 作業はISSUEを立ててからISSUE名の入ったブランチをdevelopから切って作業する  
 - 作業を始める前にdevelopをpullする  
 - 作業が完了したらdevelopにpull requestを送る  
 - margeが終わったら、作業用ブランチは削除する
 - productブランチはamazonLinux環境での動作確認用
 - masterブランチは本番用ブランチ

## ISSUEの命名ルール  
ISSUE-(番号) (作業名)
