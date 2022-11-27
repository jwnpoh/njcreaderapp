*The NJC Reader is currently still WIP.*

[![Go Report Card](https://goreportcard.com/badge/github.com/jwnpoh/njcreaderapp/backend)](https://goreportcard.com/report/github.com/jwnpoh/njcreaderapp/backend)
[![CodeFactor](https://www.codefactor.io/repository/github/jwnpoh/njcreaderapp/badge)](https://www.codefactor.io/repository/github/jwnpoh/njcreaderapp)


## Introduction
The NJC Reader is a full stack web app that is an evolution of the [NJC GP News Feed](njc-gp-newsfeed.et.r.appspot.com). This web app represents an improvement on earlier iterations in the following ways:  
- Cleaner back end code compared to previous version
- SvelteKit for a better front end experience for users
- Utilizing an actual database and the capabilities that such platforms provide, instead of Google Sheets

### New app features
#### Student/user
- A softer, more modern design and colour palette for a more pleasant user experience
- User log in/out, user profile (possibly user stats in the future)
- Personal notebook and social notes feed to check out public notes shared by other students

#### Teacher/admin
- Adding and deleting articles can now be performed on multiple articles at a go
- User management (to be implemented)
- Monitoring of teaching classes (to be implemented)

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
- [x] user authentication
- [x] user profile ui and api
- [x] user notebook
- [ ] social relations: following, likes, view public profile, etc.
- [ ] columns feed and admin ui and api
- [ ] about page
- [ ] feedback form
- [ ] user management (add user(s), edit user, delete user)
- [ ] user password reset
- [ ] profanity filter
- [ ] teacher monitoring of classes (maybe/not?)
- [ ] articles/notes/users stats

