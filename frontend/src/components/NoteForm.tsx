import Tag from "../types/Tag"
import Note from "../types/Note"
import { useRef, useState, FormEvent, useEffect } from "react"
import { Row, Col, Stack, Button, Form } from "react-bootstrap"
import { useNavigate, Link } from "react-router-dom"
import CreatableReactSelect from "react-select/creatable"

type Props = {
    onSubmit: (data: { header: string, body: string, tags: Tag[] }) => void
    AvailableTags: Tag[]
} & Partial<Note>


function NoteForm({ onSubmit, header = "", body = "", tags = [], AvailableTags }: Props) {
    const headerRef = useRef<HTMLInputElement>(null)
    const markdownRef = useRef<HTMLTextAreaElement>(null)

    const navigate = useNavigate()

    console.log(tags)
    const [selectedTags, setSelectedTags] = useState<Tag[]>(tags)

    function NoteSubmitHandler(e: FormEvent) {
        e.preventDefault()

        let data = {
            header: headerRef.current!.value,
            body: markdownRef.current!.value,
            tags: selectedTags
        }
        onSubmit(data)

        navigate("..")

    }

    return (
        <Form onSubmit={NoteSubmitHandler}>
            <Row className="align-items-center mb-3">
                <Col>
                    <h2>my neat.ly</h2>
                </Col>
                <Col xs="auto">
                    <Stack direction="horizontal" gap={2} className="justify-content-end">
                        <Button type="submit">
                            Save
                        </Button>
                        <Link to="..">
                            <Button>
                                Cancel
                            </Button>
                        </Link>
                    </Stack>
                </Col>
            </Row>
            <Row>
                <Stack gap={4}>
                    <Row>
                        <Col>
                            <Form.Group className="mb-3" controlId="header">
                                <Form.Label>Header</Form.Label>
                                <Form.Control placeholder="header" required
                                    ref={headerRef}
                                    defaultValue={header} />
                            </Form.Group>
                        </Col>

                        <Col>
                            <Form.Group className="mb-3" controlId="tags">
                                <Form.Label>Tags</Form.Label>
                                <CreatableReactSelect placeholder="tags"
                                    value={selectedTags.map(tag => {
                                        return { label: tag.label, value: Number(tag.id) }
                                    })}
                                    onChange={tags => {
                                        setSelectedTags(
                                            tags.map(tag => {
                                                return { label: tag.label, id: Number(tag.value) }
                                            })
                                        )
                                    }}
                                    onCreateOption={label => {
                                        const newTag = { id: 0, label: label }
                                        setSelectedTags(prev => [...prev, newTag])
                                    }}
                                    options={AvailableTags.map(tag => {
                                        return { label: tag.label, value: tag.id }
                                    })}
                                    isMulti />
                            </Form.Group>
                        </Col>
                        <Form.Group className="mb-3" controlId="markdown">
                            <Form.Label>Body</Form.Label>
                            <Form.Control placeholder="type here..."
                                ref={markdownRef}
                                as="textarea"
                                rows={20}
                                defaultValue={body} />
                        </Form.Group>
                    </Row>
                </Stack>
            </Row>
        </Form>
    )
}

export default NoteForm;