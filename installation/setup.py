from setuptools import setup, find_packages

setup(
    setup_requires=["setuptools_scm"],
    use_scm_version={"root": "../"},
    package_dir={"": "../pysrc"},
    name="farkle",
    packages=find_packages(),
)
