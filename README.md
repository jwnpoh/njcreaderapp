*The NJC Reader is currently still WIP.*

[![Go Report Card](https://goreportcard.com/badge/github.com/jwnpoh/njcreaderapp/backend)](https://goreportcard.com/report/github.com/jwnpoh/njcreaderapp/backend)
[![CodeFactor](https://www.codefactor.io/repository/github/jwnpoh/njcreaderapp/badge)](https://www.codefactor.io/repository/github/jwnpoh/njcreaderapp)


## Introduction
The NJC Reader is a full stack web app that is an evolution of the [NJC GP News Feed](njc-gp-newsfeed.et.r.appspot.com). This web app represents an improvement on earlier iterations in the following ways:  
- Cleaner back end code compared to previous version
- SvelteKit for a better front end experience for users
- Utilizing an actual database and the capabilities that such platforms provide, instead of Google Sheets

### Technologies/Platforms/Frameworks
#### Backend
- [Go](https://go.dev/)
- [Go-chi](https://go-chi.io/)
- [sqlx](http://jmoiron.github.io/sqlx/)
- [Planetscale](https://planetscale.com/)
- [Firebase/Firestore](https://firebase.google.com/)

#### Frontend
- [Svelte/SvelteKit](https://kit.svelte.dev/)
- [daisyUI](https://daisyui.com/)
- [Tailwind CSS](https://tailwindcss.com/)

### Progress tracking
- [x] articles feed api
- [x] articles feed ui
- [x] articles admin api
- [x] articles admin ui
- [ ] columns feed and admin ui and api
- [x] user authentication
- [x] user profile ui and api
- [x] user notebook
- [ ] social relations: following, likes, view public profile, etc.
