## Shall We Retry?

Some function calls, when they fail, are worth retrying. We had a couple of
these calls in our application so we built a small retry utility. It's small
because it still leaves some of the heavy-lifting to you. Most importantly,
interface conversion is left to you. Also it accepts only one function type, but
it's broad enough to cover you use case. All you need to do is wrapper function
that knows how to call your real function.

It does not recover from `panic`s, neither does it abort remaining trials when
specific errors occur. If you need any of these, please feel free to submit a PR
to add them.
