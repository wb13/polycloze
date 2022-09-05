"""Provides functions for checking if a target has to be rebuilt."""

from concurrent.futures import as_completed, Future, ProcessPoolExecutor
from datetime import datetime
from graphlib import TopologicalSorter
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


class TaskSummary(t.NamedTuple):
    start: datetime
    end: datetime
    name: str

    def __str__(self) -> str:
        layout = "%Y-%m-%d %H:%M:%S"
        start = self.start.strftime(layout)
        end = self.end.strftime(layout)
        return f"{start} - {end} [{self.name}]"


class WorkloadSummary(t.NamedTuple):
    tasks: list[TaskSummary]

    def __str__(self) -> str:
        return "\n".join(str(task) for task in sorted(self.tasks))


Task = t.Callable[[], t.Any]


class DependencyGraph:
    def __init__(self) -> None:
        self.sorter: "TopologicalSorter[Task]" = TopologicalSorter()
        self.tasks: set[Task] = set()

    def add(self, task: Task, *dependencies: Task) -> None:
        self.sorter.add(task, *dependencies)
        self.tasks.add(task)
        self.tasks.update(dependencies)

    def _execute_some_tasks(
        self,
        executor: ProcessPoolExecutor,
    ) -> list[TaskSummary]:
        """Execute tasks until no progress can be made.

        Returns list of completed tasks.
        """
        completed = []
        futures = []
        while self.sorter.is_active():
            sleep(0)   # Be nice to others

            for task in self.sorter.get_ready():
                def create_callback(
                    task: Task = task
                ) -> t.Callable[[Future[Task]], None]:
                    self.tasks.remove(task)
                    start = datetime.now()

                    def callback(_: Future[Task]) -> None:
                        self.sorter.done(task)
                        summary = TaskSummary(
                            start=start,
                            end=datetime.now(),
                            name=task.__name__,
                        )
                        completed.append(summary)
                    return callback

                future = executor.submit(task)
                future.add_done_callback(create_callback())
                futures.append(future)

        # Since no progress can be made, just wait for futures to complete
        for future in as_completed(futures):
            future.result()
        return completed

    def execute(self) -> "WorkloadSummary":
        """Execute all tasks in DAG.

        Returns summary of executed tasks.
        """
        completed = []

        self.sorter.prepare()
        with ProcessPoolExecutor() as executor:
            while self.tasks:
                completed.extend(self._execute_some_tasks(executor))
        return WorkloadSummary(completed)
