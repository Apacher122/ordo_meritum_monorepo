# Ordo Meritum - My Over-Engineered Job Application Co-pilot

## Looking for a job? Me too, buddy.

This is a personal tool I'm developing for myself to get through unemployment purgatory, but feel free to look around the repo. What started as a simple backend has morphed into a full-blown monorepo.

Long story short, it's an overglorified resume/cover letter proofreader that's also an application tracker and job match analyzer. The purpose of the job match analyzer is mostly to provide some extra insight on my alignment to a particular role before I fire my resume into the void.

### Previews

Soon...
<!-- ![alt text](https://github.com/Apacher122/job_hunter/blob/main/previews/Screenshot%202025-09-16%20at%2018.22.58.png "Match Score Sample")
![alt text](https://github.com/Apacher122/job_hunter/blob/main/previews/Screenshot%202025-09-16%20at%2018.23.16.png "Resume Sample")
![alt text](https://github.com/Apacher122/job_hunter/blob/main/previews/Screenshot%202025-09-16%20at%2018.24.08.png "Cover Letter Sample") -->

## Architecture - What is this thing now?

The whole setup is a collection of services that talk to each other, all living happily in this monorepo.

* `electron-frontend`: The face of the operation. A desktop app built with **Electron** and **React** so I can use this thing on my MacBook.
* `go-server`: The new brain of the operation, written in **Go**. It handles most of the core logic, talks to the database, and manages user data.
* `documents-service`: A **TypeScript** microservice that wrangles LaTeX to build the resumes and cover letters.
* **Kafka**: The nervous system. It passes messages between the Go server and the document service so they don't have to talk to each other directly.

## Getting it Running

### **Requirements**

* Docker
* Go
* Node.js & npm
* An API Key from your favorite LLM provider (Gemini, OpenAI, etc.)

### **Setting it up**

1.  **Clone the repo.** You know how to do this.

2.  **Environment Variables:** In the root directory, you'll need a `.env` file. Check out `docker-compose.yml` to see all the environment variables you need to set for the Go server and the documents service.

3.  **Install Frontend Stuff:**
    ```bash
    cd electron-frontend
    npm install
    cd ..
    ```

4.  **Light the fires:**
    From the root of the project, run the magic command:
    ```bash
    docker compose up --build
    ```
    This will build the containers for the backend services and spin everything up.

## Additional Info

Credit for the resume template goes to Claud D. Park <posquit0.bj@gmail.com>
You can view the template here: <https://github.com/posquit0/Awesome-CV>

## TODO

1.  ~~Add functionality to showcase relevant skills on resume.~~
2.  ~~Add functionality to showcase relevant projects on resume.~~
3.  ~~Rework backend.~~
4.  ~~Add cool UI stuff.~~
5.  Actually make the UI stuff cool.
6.  Figure out the logistics of opening this up for public consumption.
7.  Don't cry.