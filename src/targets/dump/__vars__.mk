# Define an array of default excluded paths
EXCLUDE_PATHS := \
    docker/ressource/** \
    package-lock.json \
    .DS_Store \
    node_modules/** \
    *.log \
    *.swp \
    *.swo \
    *.tmp \
    *.temp \
    *.bak \
    *.cache \
    .vscode/** \
    build/** \
    dist/**

# Define file extensions to exclude from content dump
CONTENT_EXCLUDE_EXT := \
    ico png jpg jpeg gif svg webp bmp tiff \
    wav mp3 mp4 ogg \
    psd ai eps raw heic \
    pdf doc docx xls xlsx ppt pptx \
    rtf odt ods odp \
    csv json xml yml yaml \
    txt md