# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information

project = 'mock'
copyright = '2024, Dhuan Oliveira'
author = 'Dhuan Oliveira'

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration

extensions = ['sphinx_design']

templates_path = ['_templates']
exclude_patterns = []



# -- Options for HTML output -------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#options-for-html-output

html_theme = 'furo'
html_static_path = ['_static']

html_css_files = [
  'css/custom.css',
]

html_sidebars = {
   '**': [
      'sidebar/brand.html',
      'sidebar/search.html',
      'sidebar/scroll-start.html',
      'sidebar/navigation.html',
      'sidebar/ethical-ads.html',
      'sidebar/scroll-end.html',
      'sidebar/variant-selector.html',
      'ga.html',
   ],
}

