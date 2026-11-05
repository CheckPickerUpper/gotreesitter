actual open class TestBase actual constructor(): OrderedExecutionTestBase(), ErrorCatching by ErrorCatching.Impl() {
    actual fun println(message: Any?) {
    }
    actual fun runTest(
        unhandled: List<(Throwable) -> Boolean>,
        block: suspend CoroutineScope.() -> Unit
    ) {
    }
}
