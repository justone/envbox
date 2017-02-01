# envbox - Secure environment variables via NaCl secretbox

Have you ever felt squirely about this?

```
$ export GITHUB_TOKEN=abcabcabcabcabcabc
$ some-command --that needs --github authentication
```

Now that environment variable is in your shell's history, not to mention that
it's exposed to every command that you run.

Wouldn't it be nice if you could store environment variables like
`GITHUB_TOKEN` encrypted and expose them only to the commands that need them?
Well, that's what envbox tries to do.

# Setup

## 1. Install

Grabbing one of the [releases](https://github.com/justone/envbox/releases) or use [holen](https://github.com/holen-app/holen).

## 2. Set key

Generate and set a key:

```
$ envbox key generate --set
```

## 3. Store an environment variable

```
$ envbox add -n GITHUB_TOKEN
value: abcabcabcabcabc
$ envbox ls
GITHUB_TOKEN=abcabcabcabcabc
```

## 4. Run commands that need those environment variables

Envbox will add the variable to the environment and then run the command.

```
$ envbox run -e GITHUB_TOKEN -- some-command --that needs --github authentication
```

For ease of use, set up an alias.

```
$ alias some-command="envbox run -e GITHUB_TOKEN -- some-command"
$ some-command --that needs --github authentication
```
