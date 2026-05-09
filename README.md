# ingit
An integrated TUI git manager that can be added to any IDE that has support for a terminal, yes i know it's like magit, but i promise there are some interesting stuff here as well.

## What's 'Ingit'
Well it stands for integrated git, cli is really cool and useful when we have to resolve something complex, but typing git this git that git all the time for people that like to commit more often can become really annoying.
Not only that but working in a mono repo enviroment is also kind of pain without external tools, modern day software development has advanced quite a lot compared to when the git cli client was created and in todays workflow it almost feels like a setback.
This tool aims to solve that, with few tricks up it's sleeve for those willing to master the dark arts terminal based applications, i will show you how to mend files with ease, alter the very structure of your source and divide those 10-20 file changes into a nicely organized execution plan

###
Features

 - Multi-Repository change overview | ideal for a mono repository projects, if you have to push changes in 3 different projects there is no need to switch editors or terminals you can simply manage them one simple interface
 - Quick commits | want to push something quickly into your current branch, no problem one key press away
 - Context aware? Well if there are no more than one repo you wont even see the project manager so it will just straight into the repository that you have
 - Easy to integrate into any IDE that has support for terminals | Amazing for i3 based systems
 - Push planner, lets you organize big changes into multiple branches, if you didn't plan ahead before committing you can create branches on the fly mark and move files and excute them all at once in a simple non destructive UI, repository and file states are not mutated till a plan is executed, so your files remain safe.


#### Integratins
If you find it interesting and are willing to try it out follow the build steps to compile the software, as for integrations, the tool comes with ready made examples for vscode and i3, because it's what i personally use but i am happy to provide more if you open a PR for your editor, wont go on a witchunt, but will most definantly make as simple as possible to integrate into the most commonly used IDEs, as long as there is the demand for it.

##### Building from source

 - First get the latest GOLANG version if you don't have it
 - navigate to the cloned folder
 - go mod tidy
 - go build -o ingit .

You can either add it to profile , zshrc, system env, or integrate it inside codium, vscoe, emacs, or figure a way to set it up for your own IDE of choice if there is no option here

I am open to PRs and suggestion for new features that you might think are useful and i am going to add them if they benefit all users.
