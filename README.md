# Clever Initiative Challenge

The Clever-Initiative is a team of the Technology Group @ Ubisoft Montreal. Our goal is to discover, improve and promote state-of-the-art software engineering practices that aim to ease the production of quality. By augmenting the quality of our products, we mechanically improve the productivity of our teams as they can focus creating features rather than fixing bugs.

The Clever-Initiative is the team behind the [Submit-Assitant](https://montreal.ubisoft.com/en/ubisoft-la-forge-presents-the-commit-assistant/) (a.k.a  Clever-Commit) that received some press coverage recently: [Wired](http://www.wired.co.uk/article/ubisoft-commit-assist-ai), [MIT](https://www.technologyreview.com/the-download/610416/ai-can-help-spot-coding-mistakes-before-they-happen/), and [more](https://www.google.ca/search?q=commit+assistant+ubisoft). The scientific foundation behind our work have been accepted for publication to [MSR'18](https://montreal.ubisoft.com/en/clever-combining-code-metrics-with-clone-detection-for-just-in-time-fault-prevention-and-resolution-in-large-industrial-projects-2/), [CPPCON'18](https://www.youtube.com/watch?v=QDvic0QNtOY).

We are also collaboration with Mozilla. More information [here](https://news.ubisoft.com/en-us/article/344442/Ubisoft-and-Mozilla-Partner-To-Develop-AI-Coding-Tools) and [here](https://blog.mozilla.org/futurereleases/2019/02/12/making-the-building-of-firefox-faster-for-you-with-clever-commit-from-ubisoft/)

We are currently looking for trainees to join us (Winter'19). The length and start date of the internship will be discussed on a per applicant basis.

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

To enter the challenge:

- Fork this repository
- Implement your solution
- Open a pull request with your solution. In the body of the pull request, you can explain the choices you made if necessary.

You can alter the data structure, add files, remove files ... you can even start from scratch in another language if you feel like it.
However, note that we do use golang internally.

If you don't feel comfortable sharing your pull-request with the world (pull-request are public) you can invite me (@mathieunls for github, bitbucket and gitlab) and Florent Jousset (@heeko for github, bitbucket and gitlab) to a private repo of yours. Don't bother to send code by email, we won't read it.

### Part 1
The first part of the challenge consists in parsing a few diffs--in the most efficient way possible (speed, maintainability, evolvability, ...)--and compute the following statistics:

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

### Part 2
The second part of the challenge consists in manipuling [ASTs](https://en.wikipedia.org/wiki/Abstract_syntax_tree). The problem consist of parsing an AST and return all the declared variables in the given format _{variable-type}{variable}_:

```
{int}{myInt}
{string}{myInt}
{Foo}{myFooObject}
...
```

You can find the AST in a JSON format [here](ast/astChallenge.json) and a visualization of it [here](ast/astChallengeViz.svg). This AST was generated for this piece of **C#** code:
```C#
void Main()
{
	//while bubble sort
	var mas = new int[] { 5, 8, 0, 5, 2, 3 };
	bool t = true;
	while (t)
	{
		t = false;
		for (int i = 0; i < mas.Length-1; i++)
		{
			if (mas[i] > mas[i + 1])
			{
				int temp = mas[i];
				mas[i] = mas[i + 1];
				mas[i + 1] = temp;
				t = true;
			}
		}
	}
}
```

## Part 3
The third part of the challenge investigates sequences. In the `seq` folder, you will find two csv files. The first file `sample.csv` contains ~18k events that are classified into two classes: `0` and `1`.
All events are uniquely identified by their `id` and occurs at a precise `timestamp`.
In addition to the `id`, the `timestamp` and the `class` each event is further categories using 30 metrics (`f1`, ..., `f30`).

Another file, named `res.csv` contains ids of the ressources used by the events describeb in `sample.csv`.

In this part, we are interested in predicting, in advance, what will be the class of the next 1-,5- and 10- events, given the information known at a precise time. For example, using the all relevent information we know after event `829454` (`f1` to `f30` and ressources information) can you predict the class of event `829455`, `829455@829459` and `829455@829464`.

You are free to use the technics / languages of your choice.


## Permanent Positions

We are also looking for permanent members to join our team. If you are interested, mail our human resource [contact](mailto:alison.laplante-rayworth@ubisoft.com?subject=Clever-Initiative) with your resume. You can submit your pull request for the challenge. However, you'll be subjected to an in-depth (much harder) coding test. This one has been conceived for students only and it might not be worth your time to take it ;).

- [x] Software Engineer / ML Dev (python, go, ml, sql, ...)
- [x] Software Engineer / Backend Dev (go, c, cfg, ast, k8s, redis, ...)
- [ ] Software Engineer / Tool devs (csharp, python, cfg, ast, ...)
