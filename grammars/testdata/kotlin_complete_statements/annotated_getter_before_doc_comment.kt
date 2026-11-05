internal val A.b: C
    @OptIn(I::class)
    get() = g(this)

/** Doc. */
fun f() {}
