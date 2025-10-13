
# Looking for a job? Me too buddy

This is a personal tool I'm developing for myself to help me get through unemployment purgatory, but feel free to look around the repo.

**As of 6/12, this is only the backend. Frontend interfaces are being developed separately so I can access this tool from my iPhone or macbook when I'm away from home.**

## For the curious

Long story short, it's an overglorified resume/coverletter proofreader that's also an application tracker and job match analyzer. The purpose of the job match analyzer is moreso to provide some extra insight on my alignment to a particular role.

## Documentation

For full documentation, visit the [Wiki](https://github.com/<your-username>/<your-repo>/wiki).

### Previews

![alt text](https://github.com/Apacher122/job_hunter/blob/main/previews/Screenshot%202025-09-16%20at%2018.23.16.png "Resume Sample")


![alt text](https://github.com/Apacher122/job_hunter/blob/main/previews/Screenshot%202025-09-16%20at%2018.24.08.png "Cover Letter Sample")


![alt text](https://github.com/Apacher122/job_hunter/blob/main/previews/Screenshot%202025-09-16%20at%2018.22.58.png "Match Score Sample")


## Important Note

**I'm currently working on making this easily accessible to others**  
My front-end is being developed as an electron app. 
I'm currently working on refactoring my server-side code to support other users outside of just myself. This means refactoring my database architecture and focusing more on data-privacy and security.

As of 9/23/2025, the road map is as follows:
1. Rework backend.
2. Create a clean electron frontend.
3. Figure out the logistics of opening this up for public consumption.

### **Requirements** (subject to change)

- Docker
- Google Cloud credentials with Sheets API enabled
- OpenAI API Key

### **Setting it up**

1. In `/root/`:
    - .env file (example provided of what info to put in there)
    - Add OpenAI, Google Sheets, and postgresql environment variables to docker-compose.yml
    - Add credentials.json from Google Cloud under `/google_config/`
2. In `/root/data/`
    - Add  
    - Add position, title, url, and description of a job you're applying for to `jobPosting.txt`
    - Add any cover letter corrections you need Open AI to be aware of if it made mistakes to `/corrections/coverLetterCorrections.txt`
    - Add any resume mistakes you need Open AI to be aware of if it made mistakes to `/corrections/mistakesMade.txt`
    - Add one or more examples of your OWN ORIGINAL writing to `user_info/my_writing/`
        - **NOTE:** It should be able to read in text, .docx,. or pdfs.
    - Add your current resume in a json format to `user_info/resume.json`.
    - Add some additional information about yourself to `user_info/aboutMe.txt`.
        - This can be a novel, autobiography, or even another resume.
3. In your editor's terminal (or whatever terminal you're using) run the following in order:

    ```bash
    docker compose build
    docker compose up
    ```

4. When the server receives an API call, it creates the following files/folders under `/root/output/`:
    - `/cover_letters/` : Contains all generated cover letter drafts, each tailored to a specific job application.
    - `/guiding_answers/` : Stores extracted or generated guiding insights to help you craft personalized applications (e.g., relevant keywords, company context, etc.).
    - `/match_summaries/` : Contains summaries analyzing how well your current resume matches a job posting, often used to guide resume revisions.
    - `/resumes/` : Holds all revised resumes, typically one per company or job application.
    - `/change-summary.md` : A human-readable markdown report summarizing all modifications made to your resume for a given job.
    - `*.log, *.aux, etc. (LaTeX build artifacts)`: Used internally by the LaTeX compiler for formatting your final documents—usually safe to ignore or delete after generation.

## Additional Info

Credit for the resume template goes to Claud D. Park <posquit0.bj@gmail.com>  
You can view the template here: <https://github.com/posquit0/Awesome-CV>

## TODO

1. ~~Add functionality to showcase relevant skills on resume.~~
2. ~~Add functionality to showcase relevant projects on resume.~~
3. ~~Add functionality to showcase relevant courses taken on resume.~~
4. Add cool UI stuff.
5. Don't cry.
