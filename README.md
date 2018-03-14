# Clever Initiative Challenge

The Clever-Initiative is a team of the Technology Group @ Ubisoft Montreal. Our goal is to discover, improve and promote state-of-the-art software engineering practices that aim to ease the production of quality. By augmenting the quality of our products, we mechanically improve the productivity of our teams as they can focus creating features rather than fixing bugs.

The Clever-Initiative is the team behind the Submit-Assitant (a.k.a  Clever-Commit) that received some press coverage recently: [Wired](http://www.wired.co.uk/article/ubisoft-commit-assist-ai), [MIT](https://www.technologyreview.com/the-download/610416/ai-can-help-spot-coding-mistakes-before-they-happen/), and [more](https://www.google.ca/search?q=commit+assistant+ubisoft).

We are currently looking for trainees to join us (Summer'18/Fall'18). The length and start date of the internship will be discussed on a per applicant basis.

## Trainees

A trainee applicant must:

- Be engaged in a computer science (or related) university program.
- Be able to work in Canada legally.
- Be willing to come to Montreal.
- Be able to read, understand and implement scientific papers.
- Know:
    - versionning systems (git, perforce, ...)
    - c/c++/csharp or java
- Know or be willing to learn:
    - golang
    - docker
    - sql
    - angular

## The Challenge

The challenge for trainee applicant consists in parsing a few diffs--in the most efficient way possible (speed, maintainability, evolvability, ...)--and compute the following statistics:

- list of files in the diffs
- number of regions
- number of lines added
- number of lines deleted
- list of function calls seen in the diffs and their number of calls

All these stats are to be computed globally (i.e. for all the diffs combined).

In the main.go file; you'll find the `compute` method that needs to be implemented.

```golang
//compute parses the git diffs in ./diffs and returns
//a result struct that contains all the relevant information
//about these diffs
//	list of files in the diffs
//	number of regions
//	number of line added
//	number of line deleted
//	list of function calls seen in the diffs and their number of calls
func compute() *result {

	return nil
}
```

To enter the challenge:

- Fork this repository
- Implement your solution
- Open a pull request with your solution. In the body of the pull request, you can explain the choices you made if necessary.

You can alter the data structure, add files, remove files ... you can even start from scratch in another language if you feel like it.
However, note that we do use golang internally.

If you don't feel comfortable sharing your pull-request with the world (pull-request are public) you can invite me (@mathieunls for github, bitbucket and gitlab) and Florent Jousset (@heeko for github, bitbucket and gitlab) to a private repo of yours. Don't bother to send code by email, we won't read it.

## Permanent Positions

We are also looking for permanent members to join our team. If you are interested, mail our human resource [contact](mailto:alison.laplante-rayworth@ubisoft.com?subject=Clever-Initiative) with your resume. You can submit your pull request for the challenge. However, you'll be subjected to an in-depth (much harder) coding test. This one has been conceived for students only and it might not be worth your time to take it ;).

- [ ] Software Engineer
- [ ] Software Engineer
- [ ] Software Engineer / Data Engineer
