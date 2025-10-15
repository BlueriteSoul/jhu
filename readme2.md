# For posterity

I'd like to document some difficulties, that could lead to enlightenment.

## Regarding jhu -locf

The original idea was simple. Bundle Tokei with my app, as my jhu -loc approach turned out to be so naive, that it didn't work for anything but the simplest of use cases. AI gave me some options have to do this, but said embedding has performance implications, so I chose another way. I hacked together a FFI for a function, that returns a number, that corresponds with the number total lines of code, from Tokei cli. But that's not what I wanted, I wanted my embedded Tokei to print the table, as if I called tokei.

I can't tell you what did I do exactly in Tokei, as I didn't take a snapshot of the code, but I think it was fairly simple. I had to modify Cargo.toml to output .a/.so files and then I think I just added few lines of code to /src/lib.rs:

```use std::ffi::{CStr, CString};
use std::os::raw::c_char;

/// # Safety
/// The `dir` must be a valid null-terminated C string.
#[no_mangle]
pub extern "C" fn count_loc(dir: *const c_char) -> u64 {
    if dir.is_null() {
        return 0;
    }

    // Convert C string to Rust string
    let c_str = unsafe { CStr::from_ptr(dir) };
    let path = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return 0,
    };

    // Prepare Tokei
    let paths = &[path];
    let excluded: &[&str] = &["target", ".git"];
    let config = Config::default();

    let mut languages = Languages::new();
    languages.get_statistics(paths, excluded, &config);

    // Sum lines of code across all languages
    let total: u64 = languages.iter().map(|(_, lang)| lang.code as u64).sum();
    total
}
```

I think that was just about everything I did. While the actual code is lost, I did preserve the snapshot of my app, including the .a/.so files, so my app's code can be examined and ran as I ran it at head: 0022f116a4cceff79e33bbac292f0b5fe7d51f65

As I didn't know what I was doing, I proceeded to wrap Tokei's main function into a custom function (I created cli_runner.rs for this) and attempted to compile this as I did before. This turned out to be a silly idea. I still don't fully understand why, but after failing to get the desired behaviour, I asked for help. I tired boot.dev discord, but nothing helpful came from that (regular St. Peert said: "On closer inspection, what you did does seem fairly close to something that should just work"). I was a bit cheeky and opened an issue on [Tokei repo](https://github.com/XAMPPRocky/tokei/issues/1286). So far there has been no response, but I didn't expect any help from there, as this is a newbie issue and these guys are busy building real tooling. Finally, I after a few days, some people on Rust Discord clarified a few things for me and that helped me find correct direction again.

Before I share the thread, I'd like to say that the code is preserved, with a compiled libs, under the head: ddf66896b96307aea20c90d4c46947a1d32f8722 and while I did make some changes to the code, the only files affected are: Cargo.toml, src/main.rs, src/cli_runner.rs, src/lib.rs)

### Discord conversation

Jerry's Handy Utility (JHU): Embedding Tokei in Go

BlueSoul (Original Poster)

11/10/2025, 00:53

Okay, where to start? Sorry to say, but I'm not a Rust guy (Go is my weapon of choice ATM) and I guess some would say, this is not a Rust issue. So bear this in mind.

While coding, I wanted a few features (that might exist, but I don't know where) that I thought I could build. I present to you Jerry's Handy Utility (JHU): https://github.com/BlueriteSoul/jhu/

I was frustrated with not being able to easily find out how much code it takes to have a project built. GitHub breaks down language distribution in a project, and hitting the API, you can learn how many bytes are written in which language. But no LOC report.

I thought, why not just clone the repo, and count the lines. Well, that's what jhu -loc does, and on the first project it failed. You see, the project came with dependencies bundled, so it still had no idea how many LOC. Anyway, I looked it up and found tools already built for this. I thought, for educational purposes (anyone can run a command programmatically, jhu -ol has no problems with invoking an external utility) I should statically link a foreign library. And that brings me to Rust.

BlueSoul (Original Poster)

11/10/2025, 00:53

I chose Tokei: https://github.com/XAMPPRocky/tokei/ and started fiddling about. With the help of AI, I managed to do something, but it wasn't obvious what actually happened. I figured it out—I interfaced with just a function that only returned lines of code (no other figures from the Tokei report were included, just the total code column). This can be observed if you set the head to commit: 0022f116a4cceff79e33bbac292f0b5fe7d51f65 (I didn't have the tokei_reference there at the time, so you won't be able to see the Rust code for this behaviour, but the behaviour is observable, as I bundled the libtokei.a file).

Since this was unintended behaviour, I tried to fix it, but it turns out, it's not that simple. So I will attempt to describe the problem.

Tokei claims in its README: "Tokei is also a library allowing you to easily integrate it with other projects." Not being a Rust guy, I'm just not sure if this should help me or not. It seems to me that the way to call foreign code is by exposing a function. But I want the whole thing to run. I tried to hack it, so main.rs really just calls the real main.rs (I had to modify Cargo.toml, obviously, but really, I only touched src/main.rs, src/cli_runner.rs, src/lib.rs) and exposing that function, but it doesn't work.

Based on the symptoms, it seems that Tokei tries (and I guess succeeds) reading from the stdin, despite I am parsing the input in Go and passing that. Even if I hardcode --help as the only argument, I still get:

jhu -locf
Running Tokei via embedded Rust library...
error: invalid value 'cf' for '--output <output>': "cf" is not a supported serialization format

For more information, try '--help'.

11/10/2025, 00:54

No matter what I do, I can't get the intended behaviour, which is basically just "jhu -locf" becoming an alias for "tokei" and if everything goes well, I'd like to pass flags, so "jhu -locf --help" would become "tokei --help" and so on and so forth. Can anybody point me into the right direction? Is this a wrong approach? I can make it work by just invoking the command, but I'd like to successfully embed some foreign code in my project.

Any advice will be appreciated, thank you.

speedy_lex.lock().await.get(42)::<>

11/10/2025, 10:14

so why the hell are you making .a files with tokei

11/10/2025, 10:15

why not just cargo add tokei and let rust handle it

11/10/2025, 10:15

    then you dont need to generate bindings

11/10/2025, 10:15

wait nvm its in go

11/10/2025, 10:18

if you build a .a from the lib then of course you wont have the tokei cli

11/10/2025, 10:18

yeah then just tell people to install tokei and spawn a process

11/10/2025, 10:19

jhu -locf ...

11/10/2025, 10:19

just take argv[2] and onwards

11/10/2025, 10:19

and pass that to tokei

BlueSoul (Original Poster)

11/10/2025, 17:42

Can you please explain this? Why would I not be able to have Tokei CLI when compiled to .a file? (I'm happy to just be given a link with heavy reading, as long as it's on topic.)

speedy_lex.lock().await.get(42)::<>

11/10/2025, 17:42

because the .a is the library

11/10/2025, 17:42

not a cli

11/10/2025, 17:43

.a is a static library it wont have the stuff for the cli

BlueSoul (Original Poster)

11/10/2025, 17:53

I don't understand. Is this so grand computer problem, that we can never solve? Like travelling salesmen and proving that program will stop? I just want to understand this. The following clip feels relevant... Can you explain the obvious please, because I genuinely don't get it?

    [Reference to Jonathan Blow video 5:33]: "We have a great deal of complexity incurred because we have this duality between command line programs and library functions... those should just be the same thing..."

speedy_lex.lock().await.get(42)::<>

11/10/2025, 19:04

ok so people dont put their CLIs in their libraries because the CLI is a wrapper for their library that's made for the command line

bruh![moment]::<>

11/10/2025, 22:02

btw does tokei solve your problem with vendored dependencies?

11/10/2025, 22:05

also I can't tell what exactly you did with the code

11/10/2025, 22:06

you could make the tokei cli into a library I guess, but the point in the docs talks about the actual library part of the API

11/10/2025, 22:10

because I don't think there is any way to exclude vendored code in general, throughout all the different project layouts

T-Dark, Speedrun World Champion (Staff)

12/10/2025, 07:17

For what it's worth, Jonathan Blow has a lot of takes and opinions that sound great to beginners but are actually... Pretty bad. I wouldn't recommend using him for any actual advice. The particular advice you're quoting here is also just not applicable.

12/10/2025, 07:19

Yes. They should [be the same]. However, in order for this to work at all, you'd need a large chunk of commonly used programs to switch to this format.

12/10/2025, 07:20

You can't just create one program that only exposes a library and expects you to use that library from your shell because shells don't support doing that (for good reason: executable formats don't have a way to do this safely).

12/10/2025, 07:21

You also can't just create a shell that supports doing this because executable formats don't really support it, and even if they did your shell would be incompatible with nearly all software that already exists.

12/10/2025, 07:21

Finally, you can't even get executable formats to support this without OS changes.

12/10/2025, 07:22

So no, it's not an impossible problem, like the halting problem. It's worse.

12/10/2025, 07:24

Of course, debating the practicality of merging CLIs and libraries isn't really going to solve your problem... your problem here and now needs to be fixed without it.

12/10/2025, 07:26

So, back to your problem. If you build the Tokei library, it doesn't contain the Tokei CLI. If you expect to have the Tokei CLI, build the Tokei CLI.

BlueSoul (Original Poster)

12/10/2025, 07:47

When I compiled the Rust code, it compiled something like 180 deps first... That probably affects the licence, right? If I understand it correctly, all I did with Tokei code is fine, as all I needed to do, was to include a NOTICE - which I did. But what about the stuff that Tokei depends on. Do you think there is a licencing issue in my case?

12/10/2025, 07:53

Like @bruh![moment] says below: "you could make the tokei cli into a library I guess...", that's what I'm trying to do.

First follow up question is: I assume you believe I didn't compiled the library, how do I do that (and can I call the library from my Go code)?

Second follow up regards the CLI problem. If I can expose the library, why can't I expose the wrapper that's been already written (wouldn't that be DRY, and therefore good?)? Can I just compile the CLI, as @bruh![moment] suggests?

12/10/2025, 07:54

I reckon I'm trying to make the Tokei CLI a library. But the problem there is, I don't know how - and it looks like I can't even compile the library. Any pointers?

12/10/2025, 08:05

@T-Dark, Speedrun World Champion Thank you for your reply, it made a lot of sense to me. Re - Jon Blow: you say: "Finally, you can't even get executable formats to support this without OS changes." and if I'm understanding what you're saying correctly, then it all makes perfect sense, because JB talks about how he would implement a new OS.

My purpose of the exercise: enable user to run SW that they didn't install.

I could compromise my vision and do something else. I mean, people say video games do this all the time, so I don't see why can't I do it. Can't I just link it dynamically, or maybe I just have to go through the installation process. And on Linux, the result could (in theory) just be an .appimage file, so in a way it's still a convenient one file redistributable. Any thoughts on anything I just said?

'pain: 'life::<> (Ferris Contributor)

12/10/2025, 08:46

for the record, Linux totally supports .so files that are executable

speedy_lex.lock().await.get(42)::<>

12/10/2025, 09:49

why wouldnt you just also install tokei and call a tokei binary? isn't that simpler?

T-Dark, Speedrun World Champion (Staff)

12/10/2025, 10:16

That it does, but "executable libraries" isn't exactly what you'd want to implement what JB suggests. You'd want something more like a library where every function can also be called directly from a shell.

12/10/2025, 10:17

A library that also happens to be an executable is enough to bundle together library and CLI, but that is mostly just undesirable. And doesn't really solve the fact you have to duplicate effort in writing both a library and a CLI.

12/10/2025, 10:20

Define "install". To run software, the user must have that software on their computer. This is inevitable.

12/10/2025, 10:23

Are you trying to run software without installing it, or are you trying to use Tokei from non-Rust code?

12/10/2025, 10:24

Wait, you're trying to hack things together so you can call the main of Tokei. WHY WOULD YOU EVER DO THIS?

12/10/2025, 10:25

Just call the library, like a normal person. The entire reason there is a library is to let you call it instead of having to go through the main entrypoint (and thus the CLI).

12/10/2025, 10:27

I don't wanna expose all the functions one by one: The conventional way to so this is code generation. You're supposed to expose all the functions (that you need to use) one by one.

12/10/2025, 10:28

Why can't I expose the wrapper that's already been written: Because that wrapper is meant to be called from a shell. Calling it from software is a massive waste of effort and performance.

12/10/2025, 10:28

Wouldn't that be DRY, and therefore good?: DRY is not automatically good. There exists such a thing as too much DRY, which can cause downsides: incomprehensibly complex abstractions, wasted performance, and deduplicating things that actually change independently of each other.

12/10/2025, 10:29

Stop trying to DRY the universe and only DRY if the alternative is worse.

bruh![moment]::<>

12/10/2025, 13:43

Well, tbh the CLI has all the pretty printing, so you could expose cli_utils, input, etc. and use them as a library, that would be slightly faster than calling the CLI binary altogether. It could be worth it if you want to distribute it statically linked.

12/10/2025, 13:43

The default fn main parses argv and doesn't expose any way to call it with other args, so that'd also need some modification.

12/10/2025, 13:44

All that is assuming that you are able to modify Rust code, which would require learning the language at least a bit.

12/10/2025, 13:44

or nerd snipe someone to do it.

### Discord conversation aftermath

I was pleased to see that JB's idea got acknowledged, but it was no coincidence that JB was talking about libs/CLI apps in context of the OS. As one of the commenters pointed out, radically new OS would be required to make this idea a reality. The same commenter eventually looked at my code and was horrified with my attempt to wrap the main function. But equiped with help from these kind strangers, I went back to the drawing board and reviewed my options.

I decided to attempt embedding Tokei. Since I don't expect anybody to really use this tool, I am willing to see what performanace impact will this have and it seems that there is only performanace impact on first run. Althought to be fair, I have no idea what goes on under the hood. I built tokei with: `cargo install --locked tokei` and put in in the /libs folder. This way seems to work as intended.
