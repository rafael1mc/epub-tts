import { parseEpub } from '@gxl/epub-parser'

let arg = process.argv[2]

const epubObj = await parseEpub(`input/${arg}`, {
    type: 'path',
})

console.log(JSON.stringify(epubObj.sections))

// let mkSections = []

// for (let i = 0; i < epubObj.sections.length; i++) {
//     // console.log("Markdown:")
//     // console.log(epubObj.sections[i].toMarkdown())
//     // console.log("Section:")
//     // console.log(epubObj.sections[i])

//     mkSections.push(epubObj.sections[i].toMarkdown())
// }
// console.log(JSON.stringify(mkSections))

