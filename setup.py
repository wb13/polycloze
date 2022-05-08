from pathlib import Path
import re
import setuptools


def get_version() -> str:
    """Get version string from blacklist/__init__.py."""
    pattern = re.compile(r'__version__ = "(.+)"')
    init = Path(__file__).with_name("blacklist")/"__init__.py"
    result = pattern.search(init.read_text())
    assert result
    return result.groups()[0]


if __name__ == "__main__":
    setuptools.setup(
        name="blacklist",
        version=get_version(),
        author="Levi Gruspe",
        author_email="mail.levig@gmail.com",
        description='Rule-based classifiers for blacklisting "non-words"',
        long_description=Path("README.md").read_text(),
        long_description_content_type="text/markdown",
        url="https://github.com/lggruspe/polycloze-blacklist",
        packages=setuptools.find_packages(),
        classifiers=[
            "Environment :: Console",
            "Programming Language :: Python :: 3.10",
        ],
        python_requires=">=3.10",
        entry_points={
            "console_scripts": [
                "blacklist=blacklist.__main__:main",
            ]
        }
    )
