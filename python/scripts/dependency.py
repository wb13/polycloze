"""Provides functions for checking if a target has to be rebuilt."""

from concurrent.futures import as_completed, Future, ProcessPoolExecutor
from graphlib import TopologicalSorter  # pylint: disable=unused-import
from pathlib import Path
from time import sleep
import typing as t


BUILD_ALWAYS = False


def mtime(path: Path, aggregate: t.Literal["max", "min"] = "max") -> int:
    if not path.exists():
        exc = FileNotFoundError(2, "No such file or directory")
        exc.filename = str(path)
        raise exc

    if path.is_file():
        return path.stat().st_mtime_ns

    agg_fn = max if aggregate == "max" else min

    child_mtime = agg_fn(mtime(child, aggregate) for child in path.iterdir())
    return agg_fn(child_mtime, path.stat().st_mtime_ns)


def is_outdated(targets: list[Path], sources: list[Path]) -> bool:
    """Build is outdated if sources timestamp > targets timestamp.

    I.e. inputs are younger than outputs.
    Assumes all inputs exist.

    Behavior can be overridden by setting build_always to True.
    """
    if BUILD_ALWAYS:
        return True

    source_time = max(mtime(source, "max") for source in sources)
    try:
        target_time = min(mtime(target, "min") for target in targets)
        return source_time > target_time
    except FileNotFoundError:
        return True


Task = t.Callable[[], t.Any]


class DependencyGraph:
    def __init__(self) -> None:
        self.sorter: "TopologicalSorter[Task]" = TopologicalSorter()
        self.tasks: set[Task] = set()

    def add(self, task: Task, *dependencies: Task) -> None:
        self.sorter.add(task, *dependencies)
        self.tasks.add(task)
        self.tasks.update(dependencies)

    def _execute_some_tasks(self, executor: ProcessPoolExecutor) -> None:
        """Execute tasks until no progress can be made."""
        futures = []
        while self.sorter.is_active():
            sleep(0)   # Be nice to others

            for task in self.sorter.get_ready():
                def callback(
                    task: Task = task
                ) -> t.Callable[[Future[Task]], None]:
                    self.tasks.remove(task)
                    return lambda _: self.sorter.done(task)

                future = executor.submit(task)
                future.add_done_callback(callback())
                futures.append(future)

        # Since no progress can be made, just wait for futures to complete
        for future in as_completed(futures):
            future.result()

    def execute(self) -> None:
        """Execute all tasks in DAG."""
        self.sorter.prepare()
        with ProcessPoolExecutor() as executor:
            while self.tasks:
                self._execute_some_tasks(executor)
