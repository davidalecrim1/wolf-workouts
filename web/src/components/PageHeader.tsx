type PageHeaderProps = {
    title: string;
}

export const PageHeader = ({ title }: PageHeaderProps) => {
    return (
        <div className="container mx-auto my-5 px-4 py-5">
            <h1 className="text-3xl font-bold">{title}</h1>
        </div>
    )
} 