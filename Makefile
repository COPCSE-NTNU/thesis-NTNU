
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

thesis.pdf: $(SRCS)
	$(mkdir)
	pdflatex $(LATEX_FLAGS) thesis
	biber $(BIBER_FLAGS) thesis
	makeglossaries thesis
	pdflatex $(LATEX_FLAGS) thesis
	pdflatex $(LATEX_FLAGS) thesis

clean:
	-@$(RM) \
		$(wildcard thesis-gnuplottex*) \
		$(addprefix thesis,.gnuploterrors .aux .bbl .bcf .blg .lof .log .lol .lot .out .pdf .run.xml .toc .acn .glo .ist .acr .alg .glg .gls)
.PHONY: clean
