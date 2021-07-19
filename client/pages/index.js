import styles from '../styles/Index.module.css'

export default function Home() {

  const submitForm = async ()=>{
    const formData = new FormData()
    let inputFile = document.querySelector('input[type="file"]').files[0]
    let ignoredRows = document.getElementById("ignoredRows").value
    let questions = document.getElementById("questions").value
    let answers = document.getElementById("answers").value
    let comments = document.getElementById("comments").value
 
    if (inputFile === undefined) {
      console.log("File is empty")      
      return
    }

    // Append the data to the form
    formData.append("file", inputFile)
    formData.append("ignoredRows", ignoredRows)
    formData.append("questions", questions)
    formData.append("answers", answers)
    formData.append("comments", comments)

    const resp = await fetch('http://localhost:5000/parsefile',{
      method: "POST",
      body: formData
    })

    console.log(resp)
  }

  return (
    <div className={styles.container}>
      <input type='file' id='' /> <br />
      <label htmlFor='questions'>Ignored Rows:  </label>
      <input type='text' id='ignoredRows' /><br />
      Please enter number of ignored rows separated by a comma EX: 1,5,6,8<br />
      <label htmlFor='answers'>Column with Questions: </label>
      <input type='number' id='questions' /><br />
      <label htmlFor='answers'>Column Num with Answers: </label>
      <input type='number' id='answers' /><br />
      <label htmlFor='comments'>Column Num with Comments: </label>
      <input type='number' id='comments' /><br />
      <button onClick={submitForm}>Upload File</button>

    </div>
  )
}
