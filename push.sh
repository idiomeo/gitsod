git add .
echo "输入commit内容:"
read commiting
git commit -m "${commiting}"
git push gitee master
git push codeberg master
git push notabug master
git push github master