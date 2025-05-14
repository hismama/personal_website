# My Personal Website
### Deploying to GCloud
1. Run the following command to see the project your gcloud is currently pointed to:
```cmd
gcloud config get-value project
```
2. List Configured Projects:
```cmd
gcloud projects list
```
3. Minify & Bundle CSS
```cmd
cat static/css/css-reset.css \
static/css/root.css \
static/css/nav.css \
static/css/profile.css \
static/css/content.css \
static/css/footer.css \
static/css/gopher-banner.css \
static/css/gopher-toast.css \
static/css/race-timer.css \
static/css/about.css \
static/css/home.css \
static/css/resume.css \
static/css/projects.css \
static/css/game.css \
static/css/contact.css \
static/css/main.css \
> static/min/site.css
```
```cmd
minify -o static/min/site.css static/min/site.css
```
4. Purge unnecessary css
```cmd
npx purgecss --content templates/**/*.gohtml --css static/min/site.css --output .
```
5. Optimize SVG
```cmd
npx svgo -f static/img/svg
```
6. Test Locally
```cmd
dev_appserver.py app.yaml
```
7. Deploy Your Go Project:
```cmd
gcloud app deploy
```
