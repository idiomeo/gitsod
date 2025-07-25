git add .
echo "输入commit内容:"
read commiting
git commit -m "${commiting}"
git push origin gitee master
git push origin codeberg master
git push origin notabug master