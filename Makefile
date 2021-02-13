
SRCS=\
	$(wildcard appendices/*) \
	$(wildcard chapters/*) \
	$(wildcard csvtables/*) \
	$(wildcard figures/*) \
	$(wildcard listings/*) \
	$(wildcard papers/*) \
	ntnuthesis.cls \
	thesis.bib \
	glossary.tex \
	thesis.tex

LATEX_FLAGS=-shell-escape
BIBER_FLAGS=

mkdir = @mkdir -p $(@D)

build: $(SRCS)
	$(mkdir)
	pdflatex $(LATEX_FLAGS) thesis
	biber $(BIBER_FLAGS) thesis
	makeglossaries thesis
	pdflatex $(LATEX_FLAGS) thesis
	pdflatex $(LATEX_FLAGS) thesis

# generate pdf and delete extras
thesis.pdf: build clean

clean:
	-@$(RM) \
		$(wildcard thesis-gnuplottex*) \
		$(addprefix thesis,.gnuploterrors .aux .bbl .bcf .blg .lof .log .lol .lot .out .run.xml .toc .acn .glo .ist .acr .alg .glg .gls .fls .fdb_latexmk)

.PHONY: clean
